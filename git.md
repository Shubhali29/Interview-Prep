<!-- This file will contain details of some useful git commands  -->
1. git branch --set-upstream-to=origin/<branch> main  
    - This command will configure the local branch to track with remote branch. After this config, user can 
        simplify commands like "git push" and "git pull" without specifying the remote branch name.
    - It helps git to know with which remote branch to synchronize the changes.