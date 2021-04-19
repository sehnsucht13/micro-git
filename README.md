# micro-git
This project is my attempt at making a small clone of git and improving my golang skills.

**Note:** This project is not meant to be a replacement for git in any way. If you are looking for an alternative, feel free to check out [go-git](https://github.com/go-git/go-git).
## Project Goal
The main goal of this project is to create a scaled down version of `git` which is written entirely in Golang. This project will support a small subset of git's features but it will not aim to be 100% compatible with git (i.e. you will not be able to manipulate existing repostiories using this project).

## Planned Features
- [ ] Creating a bare repository
- [ ] Manipulating the staging area(`git add` and `git restore`)
- [ ] Viewing the staging area(`git status`)
- [ ] Branching
- [ ] Commiting 
- [ ] Viewing Logs
- [ ] Smaller utilities like `git cat-file` and `git hash-object`
- [ ] Cloning a repository from another local repository
- [ ] Configuration command similar to `git config`
- [ ] Ability to view an arbitrary remote repository like `git ls-remote`
- [ ] Checking out a previous commit
- [ ] Moving/deleting files (`git mv` and `git rm`)
- [ ] Support for an ignore file and `git check-ignore`
- [ ] (MAYBE) Merging
- [ ] (MAYBE) Clonging a remote repository

## Good Learning Resources
- The `Git Internals` section from [Pro Git](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain) is a fantastic start to learning more about how git works internally.
- Mary Rose Cook's [Gitlet](http://gitlet.maryrosecook.com/docs/gitlet.html) is another interesting reference as well. It uses literate programming to walk the reader through creating a basic version of git. While there is no support for the majority of flags that are available for each of git's functions, it is a fantastic overview.

