# spawn

## Goals
- Easy templating for a new chain from base

- Local-Interchain nested, easy JSON configured starts

- Chain-Chores like features (pull all patches from upstream based off the value of the current spawn instance. i.e. spawn v1.0 pulls from the v1.0 branch)

- Easily add CI/CD in line with the template repo (could just pull from this Repos CI so we can confirm all works upstream. Then wget down)

- Base for a new module into your repo (spawn module new <module-name>). Regex import into your apps SDK without any comments in the app.go

- Easily import or clone upstream modules into your repo (spawn module import <module-name>). Module name can also be a git repo (even one we do not own) that we can pull the files from directly for the user. So if we want SDK v0.50 tokenfactory, we can pull from repo X or repo Y depending on our needs. May require a unique go.mod for each module, unsure atm. Maybe we can abstract this away and handle ourselves?