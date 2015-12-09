ghuwtchr (_GitHub Un-Watcher_) unwatches GitHub repos for you.
===============

It makes sense in case a long list of repositories should be unwatched at once while avoiding using the GitHub web interface.  

### Installation:

- Requires a Go >= 1.5.1 development environment on the machine.
    - With `$GOPATH` defined in environment and `go` binary in `$PATH`.

            $ go get github.com/asemt/ghuwtchr

### Usage:

- Set your personal [GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/) in environment:

        $ export GHACCESSTOKEN=<your-access-token>

- Run the unwatch utility:

        $ $GOPATH/bin/ghuwtchr

- Now past the GitHub repository links you want to unwatch from into the shell. 
The links should be formed like this:  

        https://github.com/<owner>/<repo-name>/subscription
        https://github.com/<owner>/<repo-name>

A single dot (`.`) ends the input and starts the unwatching GitHub API calls.
