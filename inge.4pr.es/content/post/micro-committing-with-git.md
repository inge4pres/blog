---
date: "2021-04-03T17:10:26+01:00"
title: "Micro-committing with Git"
author: "inge4pres"
emoji: true
comment: true

categories:
  - tech 

tags:
  - software-engineering
  - practices
---

I have been using micro-committing for some time now, during which I have adapted the usage of this technique to my needs,
bringing it to a level that makes me more productive than ever in software development. 

Combining micro-committing with Git, while doing [TDD](https://it.wikipedia.org/wiki/Test_driven_development)
is now my favorite development experience: I like how this workflow helps to deliver changes with speed and confidence.

This is not a one-size-fits-all approach, I'm sharing what works great _for me_; I hope some parts of what follows 
will help you and your team as well.

### What is micro-committing?

Despite micro-commits are not new, there is not much literature or a formal definition of it. 
The concept is very simple: when developing or refactoring, use the SCM _at every notable step_ creating a small commit.

Sizing properly micro-commits is what takes more practice mastering: I think it's a very personal thing, one might be comfortable
with lots of very tiny commits, some other with slightly larger bundles of changes, the important factor is usability.

We want to store in a commit every relevant code that produces an increment towards a goal.

##### Why should you care?

I found that using micro-commits makes working iteratively safer: the confidence I gain in adding or changing code once
I am backed by a _"diffable history"_, is priceless.

Keeping track of changes at a smaller rate enables to move back and forth into history to do local optimizations.
Having these changes recorded with Git might not be useful _all the time_, but when they're needed they are going
to be very valuable.

### How to make the best out of it

In theory, the smaller the change bundle recorded in every commit, the better!
While keeping a long list of tiny changes locally can have several benefits, it can also turn your review process and CI
pipelines into chaos.

That's why we need to be careful in managing micro-commits not only during coding but also when preparing changes to
be submitted for peer review.

#### During development

During the development phase we want to store locally as many micro-commits as we need to ensure that we can easily
move our code back to a working revision. When the development is completed we don't want to send the micro-commits 
to the remote repository, because reviewing that history would be difficult.

With these concepts in mind, below the things to remember when mirco-committing during the development of new features.

* **never** push to the remote repository until the commit history has been reviewed

One thing that I learnt recently thanks to my [amazing colleagues](#credits) is that Git history is precious.
We don't want to pollute the tree with too many commits targeting the same feature.
Micro-commits are an aid for _local_ development: they should be squashed, rewritten and edited before sending them
into a remote repository (more on this below).

* add a micro-commit for every development iteration unit: every test assertion added, every successful implementation and 
  refactoring

We want that every change with a significant impact on the codebase can be consistently fetched.
Any IDE Undo/Redo feature does kind of the same, but a micro-commit plays better with refactoring, for example when
changing a signature or renaming a package.

We don't want to store in a micro-commit changes that do not have a "business" meaning, like renaming a variable.
At the same time we don't want to store in a micro-commit several changes that map to multiple features: doing so would 
prevent us to go back in time with fine-grained control, and thus lose the benefits or incremental history.

* write a short but meaningful commit message: we must be able to understand what we did

Optionally, include in the commit message the name of the component or sub-system under change.
When re-reading the history of commits, it will be helpful to know which commits are part of which
bundle of changes.

* create a tag when a cycle of development or refactoring is complete, but delete it before pushing 
  (in summary, never push tags used for referencing micro-commits)

This is not strictly required: I like to create a simple, non-annotated tag, with the name of the task when I declare
it completed. It helps me to point a commit to a task completion in the Git log.
Adding a tag can be replaced by details in the commit message: we need to be able to understand which commit resolved
a given task.

* when you are stuck, throw everything away and go back to last commit!

This is where this technique really shines: firing up a single command 

```bash
git reset --hard
```

will allow to re-start from what _you_ chose minutes ago as a standing point, a safe place from which a new iteration can start.
If you create tags as in the previous bullet, you can even jump multiple micro-commits back, correcting a wider design 
mistake or a faulty implementation with a newer idea, just specify the tag name you want to `reset` to.

#### Preparing for peer review

It is important to review the Git history before submitting it to the remote repository: we should group multiple
micro-commits belonging to the same task into a single commit.

The reasons for doing so:

1. peer-review done commit by commit helps reviewers understanding more than the code logic, enriching code with
   the _rationale_ and the decisions made by the author
1. reverting a single feature should be as easy as reverting a single commit (or a few of them): reducing the
   cognitive load embedded in the history helps when things go wrong, even if all tests had passed 😁

A good way to achieve a clean history, ready for being reviewed, is with an [**interactive rebase**](https://git-scm.com/docs/git-rebase)
to the target branch.
This is the most straightforward way to create a meaningful list of commits that can tell a story about the development 
being shipped.

When rebasing interactively, you are presented the list of commits you are going to add on top of the target branch in an
editor. From there it's possible to merge (`squash`, `fixup`), remove (`drop`) and even change the order of commits!

I use these 3 steps to craft the intended Git history:

1. remove tags used as references to micro-commits

In most CI systems tied to the Git tree, tags can trigger special workflows like creating artifacts for releases:
we don't want this.

If tags have been created to keep track of multiple features, I'll remove them to avoid accidentally pushing them 
and triggering such workflows.

1. `squash` or `fixup` corrections and clean-up commits to the uppermost relevant commit 
   
It is good to remove commits that were useful during local development, but do not have a tangible impact for the
reviewers: i.e. linting corrections, typos in documentation, etc... These commits should always be conflated with a previous one.

It's also possible to move commits up and down during an interactive rebase, but note that it can lead to difficulties in 
setting up the desired history: you might end up in conflicts that will counter-effect the benefits of micro-committing,
with time lost on editing during the rebase.

1. set the commit that is the union of multiple fixup/squash commits with a `reword`

In this way, Git will open the editor to revrite the commit message.
A commit that will be in the history of a repository needs to be explicit about what's being committed.
It has to contain a good high-level description of what it holds in all its parts.

Remember the famous jewelery slogan _"a diamond is forever"_? A valid statement is also __"a commit is forever"__ 
(at least it should be!) 😉.

#### A practical example


### Credits
