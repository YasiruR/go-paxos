#!/bin/bash
#num_leaders=$1
#num_replicas=$2
#cd ..
#go build
#echo "binary created"

cd ..
./run leader localhost:2022 localhost:2023,localhost:2024 localhost:2025,localhost:2026 &
./run leader localhost:2023 localhost:2022,localhost:2024 localhost:2025,localhost:2026 &
./run leader localhost:2024 localhost:2023,localhost:2022 localhost:2025,localhost:2026 &

./run replica localhost:2025 localhost:2022,localhost:2023,localhost:2024 localhost:2026 &
./run replica localhost:2026 localhost:2022,localhost:2023,localhost:2024 localhost:2025
