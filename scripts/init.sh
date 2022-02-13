#!/bin/bash
num_leaders=$1
num_replicas=$2
first_port=$3

if [ ! -f ./leaders.txt ]
  then
    :
  else
    rm leaders.txt
fi

if [ ! -f ./replicas.txt ]
  then
    :
  else
    rm replicas.txt
fi

# generating ports and leaders
leaders=()
for i in `seq 0 $((num_leaders-1))`
do
  # shellcheck disable=SC2100
  p=$((first_port+i))
  leader="localhost:${p}"
  leaders+=("${leader}")
  echo "$leader" >> leaders.txt
done

# generating ports and replicas
replicas=()
first_port=$first_port+$num_leaders
for i in $(seq 0 $((num_replicas-1)))
do
  # shellcheck disable=SC2100
  p=$((first_port+i))
  replica="localhost:${p}"
  replicas+=("$replica")
  echo "$replica" >> replicas.txt
done

# creating replica list as a string
replicas_str=""
for r in "${replicas[@]}"; do
  if [ "$replicas_str" = "" ]; then
    replicas_str+=" ${r}"
  else
    replicas_str+=",${r}"
  fi
done

# creating leader list as a string
leaders_str=""
for l in "${leaders[@]}"; do
  if [ "$leaders_str" = "" ]; then
    leaders_str+="${l}"
  else
    leaders_str+=",${l}"
  fi
done

cd ..

# starting leaders
for l in "${leaders[@]}"; do
  index=0
  command="./run leader ${l} "
  for a in "${leaders[@]}"; do
      # shellcheck disable=SC2077
      if [ "$a" = "$l" ]
      then
        continue
      else
        index=$((index+1))
        if [[ $index != 1 ]]; then
          command+=","
        fi
        command+=$a
      fi
  done
  command+=$replicas_str
#  echo "$command" >> leaders.txt
  eval "$command &"
done

# starting replicas
for r in "${replicas[@]}"; do
  index=0
  command="./run replica ${r} ${leaders_str} "
  for a in "${replicas[@]}"; do
      # shellcheck disable=SC2077
      if [ "$a" = "$r" ]
      then
        continue
      else
        index=$((index+1))
        if [[ $index != 1 ]]; then
          command+=","
        fi
        command+=$a
      fi
  done
  eval "$command &"
#  echo "$command" >> replicas.txt
done
