# AWS IAM policy to let users manage their own MFA

If you&#8217;re an AWS administrator you know that managing web console security is pretty tough unless you know what you want and you know what you&#8217;re doing. So if what you want is let each AWS user manage their own MFA device configuration without you and force them to have MFA active to use the web console, here is your solution.

**TL;DR**

  * Create one or more groups with your web users
  * Create a new policy using <a href="https://github.com/inge4pres/blog/blob/master/aws-iam-policy-to-let-users-manage-their-own-mfa/mfa-for-web-users" target="_blank">this JSON</a>
  * Attach the policy to the group(s)

**How does it work?**
  
The policy has this logic:

  * Allow basic operations on IAM without having MFA set up
  * Allow setup and management of MFA for own user in IAM (create, delete, resync device)
  * Deny every action on every resource when MFA is not setup
  * Allow user to access IAM without MFA &#8211; this is necessary to sub-segment the previous rule

The magic lies in the use of ARN policy variables which is a poorly documented feature of IAM. Notice how in some case the statement makes use of `${aws:username}` to confine the action executed on the only user receiving the policy grants.

This IAM policy blocks every serice usage when MFA is not setup, and in conjunction with default IAM behavior will deny access on every action if not explicitly given. You should combine this &#8220;base&#8221; policy with other group/service oriented policies to confine web users on certain functionalities. For example if you want a set of users self-managing their own MFA and access the EC2 service only after having setup MFA, you should execute the following after having setup the IAMUsersMFAManagement policy.

<pre class="theme:sublime-text lang:batch decode:true" title="AWS CLI create EC2 web user">aws iam create-group --group-name ec2webgroup
aws iam create-user --user-name ec2webuser
aws iam add-user-to-group --group-name ec2webgroup --user-name ec2webuser
aws iam attach-group-policy --policy-arn arn:aws:iam::AWSACCOUNTID:policy/IAMUsersMFAManagement --group-name ec2webgroup
aws iam attach-group-policy --policy-arn arn:aws:iam::aws:policy/AmazonEC2FullAccess --group-name ec2webgroup



</pre>

``Reference: [AWS IAM variables documentation][1]

 [1]: http://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_variables.html

