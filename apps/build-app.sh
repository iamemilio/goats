#!/bin/bash
set -ex

if command -v docker &> /dev/null
then
    runtime=docker
elif command -v podman &> /dev/null
then
    runtime=podman
else
    echo "unsupported container runtime: must have either docker or podman"
    exit 1
fi

show_help() {
    tee <<EOF
build-app will build your app image for you and push it up to a container registry of your choice.

Usage: $0 [options] <agent> <tag> <application>

Examples:
  build from a remote git repo: $0 -a Go -t quay.io/example_user/app_repo:latest example

tag:
  -t, --appTag     <registry/repo/image:tag>      this tag will be applied to your application image when built and pushed

agent:
  -a, --agent      <agent directory>              the agent path that the application is under  

Options:
  -h, --help                                      see this help message
  -r, --remoteBranch <remote/branch:agent>        checks out a language agent from a remote branch into the container directory
                                                    example: -r newrelic/develop:go-agent
EOF
}

args=$#
i=1
opts=$(($args-1))

if [ "$opts" -eq "2" ]
then
  echo "invalid usage: $0 $1 $2"
  echo "for help use the command $0 --help"
  exit 1
fi

if [ "$opts" -ge "3" ]
then
  while [ $i -le $opts ] 
  do
      case "$1" in
        "-r"|"--remote-branch")
          remoteBranch=$2
          shift 1
          i=$((i + 1));
          ;;
        "-t"|"--tag")
          tag=$2
          shift 1
          i=$((i + 1));
          ;;
        "-a"|"--agent")
          agent=$2
          shift 1
          i=$((i + 1));
          ;;
        "-h"|"--help")
          show_help
          exit 0
          ;;
      esac
      i=$((i + 1));
      shift 1;
  done
fi

case "$1" in
  "-h"|"--help")
    show_help
    exit 0
    ;;
  *)
    APPLICATION=$1
    ;;
esac

if [ -z "$APPLICATION" ]
then
    echo "an application must be provided"
    echo "for help use the command $0 --help"
    exit 1
fi

if [ -z "$tag" ]
then
    echo "a tag must be provided"
    echo "for help use the command $0 --help"
    exit 1
fi

if [ -z "$agent" ]
then
    echo "an agent must be provided"
    echo "for help use the command $0 --help"
    exit 1
fi

cd ${agent}
cd ${APPLICATION}

if [ ! -z "$remoteBranch" ]
then
    arrIN=(${remoteBranch//\// })
    remote=${arrIN[0]}    
    branch=${arrIN[1]}
    arrIN2=(${remoteBranch//\:/ })
    agent=${arrIN2[1]}

    if [ ! -d "go-agent" ]
    then
        echo "checking out the the $agent"
        git clone https://github.com/newrelic/${agent}.git
    fi
    echo "checking out $remoteBranch"
    cd go-agent
    git remote add ${remote} "http://github.com/$remote/$agent.git" | true
    git fetch ${remote} ${branch}
    git checkout --track ${remote}/${branch} | true
    git pull
    cd ..
fi

cd app

$runtime build -t $tag .
$runtime push $tag
