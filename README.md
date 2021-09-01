# Chore
Chore is a tool for executing common tasks you run on your remote servers. You can easily setup tasks for deployment, commands, and more.

# Installation
Download release binaries from [here]().

## Initialize Example Config
Below command will generate example taskfile.yaml.
```sh
chore init
```

# Writing Tasks
## Defining Tasks
Tasks are the basic building block of Chore. Tasks define the shell commands that should execute on your remote servers when the task is invoked. For example, you might define a task that pull latest changes from git provider `git pull origin master`.

All of your tasks should be defined in an `taskfile.yaml` file. Here's an example to get you started:
```yaml
servers:
  - user@192.168.1.1
  - user@192.168.1.2

tasks:
  - name: deploy
    run: 
      - cd /path/to/site
      - git pull origin main
```

## Server Options
By default server will use default options for port and private key file path, but you're specify different options, see exmaple below.
```yaml
#...
servers:
  - user@192.168.1.1@2222
  - user@192.168.1.2@2222@/home/user/.ssh/vps_id_rsa
#...
```

## Local Tasks
You can force a script to run on your local computer by specifying the server's as localhost or 127.0.0.1
```yaml
servers:
  - localhost # or 127.0.0.1
```

## Variables
You can set common vairables which you can re-use in multiple tasks.

```yaml
# ....
vars: 
  branch: main

tasks:
  - name: deploy
    run: 
      - cd /path/to/site
      - git pull origin $branch
# ...
```

## Environment Variables
You can set the environment variables, all these env values will be exported before running task commands.
```yaml
# ....
tasks:
  - name: deploy
    env:
        FOO: BAR
        BAR: FOO
    run: 
      - cd /path/to/site
      - git pull origin main
# ...
```

## Buckets
Buckets group a set of tasks under a single, convenient name. For instance, a deploy bucket may run the update-code and install-dependencies tasks by listing the task names within its definition:

```yaml
tasks:
  - name: pull-code
    run: 
      - cd /var/www/site
      - git pull origin main

  - name: install-dep
    run:
      - cd /var/www/site
      - yarn install

buckets: 
  - name: deploy
    tasks:
      - pull-code
      - install-dep
```

# Running Tasks
To run a task or bucket that is defined in your application's taksfile.yaml file, execute Chore's run command, passing the name of the task or bucket you would like to execute. Chore will execute the task and display the output from your remote servers as the task is running:
```sh
chore run deploy
```
## Parallel Execution
By default, tasks will be executed on each server serially. In other words, a task will finish running on the first server before proceeding to execute on the second server. If you would like to run a task across multiple servers in parallel, add the parallel option to your task declaration:
```sh
chore run deploy --parallel
```

## Usage Reference
```sh
chore run --help

Run the tasks defined in Taskfile file.

Usage:
  chore run [flags]

Flags:
      --bucket            Run the bucket of tasks
      --continue          Continue running even if a task fails
      --dry-run           Dump Bash script for inspection
      --filename string   The name of the Commet file (default "taskfile.yaml")
  -h, --help              help for run
      --parallel          Run task concurrently on servers
      --path string       The path to the Commet.yaml file
```
## LICENSE
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.