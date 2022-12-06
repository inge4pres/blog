# gRPC Traffic Mirroring With Ingress-Nginx on K8s


In a [previous post](https://inge.4pr.es/post/grpc-traffic-mirroring-using-nginx/) we saw an NGINX configuration to allow gRPC traffic mirroring.

Is the same technique applicable on Kubernetes?
Yes! Using the [ingress-nginx](https://github.com/kubernetes/ingress-nginx) ingress controller!

### Traffic mirroring

Use the following configurations snippets in the [ingress-nginx configMap](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/)
and in the Ingress manifest to mirror **all traffic** to a separate gRPC server.

#### ConfigMap

Replace `grpc-backend.company.net` and `grpc-mirror.company.net` with the original and mirror endpoint, respectively.

```yaml
    http-snippet: |
      server {
        listen 127.0.0.1:9443 ssl http2;
        server_name grpc-backend.company.net;
        grpc_ssl_protocols TLSv1.3;
        ssl_certificate_by_lua_block {
          certificate.call()
        }
        location / {
          grpc_pass grpcs://grpc-mirror.company.net:443;
        }
      }
```

#### Ingress

Add the following annotations in the Ingress manifest defining the endpoint of your gRPC service.

```yaml
    nginx.ingress.kubernetes.io/backend-protocol: GRPC
    nginx.ingress.kubernetes.io/configuration-snippet: |
      mirror /mirror;
    nginx.ingress.kubernetes.io/server-snippet: |
      location = /mirror {
        internal;
        proxy_set_header X-Mirrored-From $http_host;
        proxy_pass https://127.0.0.1:9443$request_uri;
      }
```

#### NGINX configuration details

Compared to the [non-K8s version](grpc-traffic-mirroring-using-nginx.md#the-solution) we have some differences:

In the main `nginx.conf` file, applied through the configMap, we have a weird section.
    
    ssl_certificate_by_lua_block {
      certificate.call()
    }
   
This is a something I discovered while looking in the ingress-nginx source: it's a helper used to load the right TLS certificate
which is impossible to do otherwise, because TLS certificates are stored in Kubernetes secrets, instead of normal files.
This replaces the TLS certificate loading directives.

The rest is unchanged.

The Ingress resource manifest contains the annotation to configure a gRPC backend `nginx.ingress.kubernetes.io/backend-protocol: GRPC`
and has 2 important settings.

The first snippet

    nginx.ingress.kubernetes.io/configuration-snippet: |
      mirror /mirror;

adds the mirroring directive to the virtual server location, to copy the gRPC traffic to the `/mirror` internal location.

The second snippet

    nginx.ingress.kubernetes.io/server-snippet: |
      location = /mirror {
        internal;
        proxy_set_header X-Mirrored-From $http_host;
        proxy_pass https://127.0.0.1:9443$request_uri;
     }

creates the internal location that will proxy the traffic to the additional server created in the configMap above.

That's it for copying _all_ traffic from an ingress to a separate server!
But what if we'd like to only mirror _a portion_ of the traffic?

At the end of a previous post I left as a homework for the readers to discover how to copy only a percentage of traffic.
Read on to see how to achieve it.

### Bonus: mirror a part of traffic

NGINX has a [`split_clients` module](https://nginx.org/en/docs/http/ngx_http_split_clients_module.html) that is capable
of setting a variable based on the _distribution of an input_. The variable can be used in virtual servers to apply 
conditional configurations.

The syntax to configure the module is

    split_clients <input string> <variable> {
      5% something;
      10% nothing;
      * "";
    }

with the value of `<variable>` being set based on the hash of `<input string>`: this can be anything that NGINX assigns
when processing a request.

The important detail to understand of the above configurations, is how to choose the input string: the percentage defines
the portion of hash values that will yield in `<variable>` the value to its right.

Let's have a look at the configs.

NGINX configMap:

```yaml
    http-snippet: |
      split_clients "${remote_addr}mirror${request_uri}" $mirror_backend {
        10% 1;
        * "";
      }
      server {
        listen 127.0.0.1:9443 ssl http2;
        server_name original.domain.com;
        grpc_ssl_protocols TLSv1.3;
        ssl_certificate_by_lua_block {
          certificate.call()
        }
        location / {
          grpc_pass grpcs://destination.domain.com:443;
        }
      }
```

Ingress manifest:

```yaml
    nginx.ingress.kubernetes.io/backend-protocol: GRPC
    nginx.ingress.kubernetes.io/configuration-snippet: |
      mirror /mirror;
    nginx.ingress.kubernetes.io/server-snippet: |
      location = /mirror {
        internal;
        if ($mirror_backend = "") { 
          return 200; 
        }
        proxy_set_header X-Mirrored-From $http_host;
        proxy_pass https://127.0.0.1:9443$request_uri;
      }
```

The additions to the previous config: with `split_clients` in the main nginx.conf file we set a variable `$mirror_backend`
with a non-empty string when the hash of `"${remote_addr}mirror${request_uri}"` falls in the _first_ 10% of all possible
hash values. 
The `if` added in the Ingress manifest will only proxy traffic when the `$mirror_backend` variable is not empty.

The hash value can be from 0 to 4294967295 (NGINX uses [MurMurHash2](https://en.wikipedia.org/wiki/MurmurHash#MurmurHash2),
returning a 32-bit integer); the percentages written in the configuration create segments of the whole hash space in a 
contiguous manner, starting from 0.

In the example above where we defined `5%`, `10%` and `*`, you will have 3 ranges of possible values for `<variable>`:
* `5%` -> hash values from 0 to 214748364
* `10%` -> hash values from 214748365 to 429496728
* `*` -> all remaining hash values from 429496729 to 4294967295

Therefore, the probability of getting each value is not _exactly_ the same of the percentage configured,
because the distribution of hashes tightly depends on the input.

For example, if you have most of your traffic from a few `$remote_address`es, you don't want to set it as input alone.
The more you have a sparse distribution of values in the input variable, the better the filter will work. 

This is why in the example configuration, I added `$request_uri` to the input, concatenating with a constant string 
`mirror`: this highly increases the entropy of the hashes, making the percentage more reliable.

An important property of hashing the input is that it's deterministic, so if you want to mirror _exactly_ for a subset
of requests, you can do it: define the percentages to include only the portion of hashes that you want to be copied.

For example, if you want to mirror only traffic for a certain URI, as in:

    $request_uri        = `/bank.Service/askForTransactions`
    murmurhash2         = 1227040391
    range percentage    = 28.57%

then the `split_clients` config would be

    split_clients $request_uri $mirror_backend {
      28.5699% "";
      28.57% 1;
      * "";
    }

### Drawbacks and limitations

Copying traffic using this technique is simple and effective, but it has a cost: we have a number of TCP connections
that are dedicated to serving the cloned traffic, even if going through the loopback interface.

You will notice from the ingress-nginx controller metrics that enabling the virtual server via the configMap does not
create more connections immediately, but as soon as you configure an Ingress to mirror using the new server, there will
be an increase in the average open connections.

This is expected, because of the non-native way we are doing mirroring for gRPC traffic.

The same applies to memory and CPU usage: handling more connections, and decrypting and re-encrypting every gRPC call
will come with a resource overhead.

One more thing to note: this technique _might_ be working with gRPC streams, but I was only able to test it with unary
RPCs.


### Credit

Thanks again to Joni ([@mejofi](https://twitter.com/mejofi)) for helping me find the original gRPC traffic mirroring configuration.

The partial mirroring addition is taken from this nice blog post by Alex Dzyoba 
https://alex.dzyoba.com/blog/nginx-mirror/

