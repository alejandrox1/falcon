#!/usr/bin/env bash
clear
set -e

printf "\n"
qsub -N P97f7c4234bfa71 -pe smp 8 -q your_queue -V -cwd -o stdout -e stderr -S /bin/bash /ecoli_test/mypwatcher/wrappers/run-P97f7c4234bfa71.bash
