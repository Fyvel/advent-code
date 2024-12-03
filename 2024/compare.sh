#!/bin/bash

if [ $# -eq 0 ]; then
	echo "Usage: $0 <target_folder> [number_of_runs]"
	exit 1
fi

TARGET_DIR="$1"
RUNS=${2:-5} # Default to 5 runs if not specified

if [ ! -d "$TARGET_DIR" ]; then
	echo "Error: Directory '$TARGET_DIR' does not exist"
	exit 1
fi

cd "$TARGET_DIR" || exit 1

# Go build
if [ -f "main.go" ]; then
	echo "Building Go executables..."
	go build -o advent_code_2024_go main.go
fi

# Rust build
if [ -f "main.rs" ]; then
	echo "Building Rust executable..."
	mkdir -p "src"
	cp "main.rs" "src/main.rs"

	# Generate Cargo.toml
	cat <<EOF >"Cargo.toml"
[package]
name = "advent_code_2024_rust"
version = "0.1.0"
edition = "2021"

[dependencies]
tokio = { version = "1.28", features = ["full"] }
EOF

	cargo build --release
	rm -rf "src" "Cargo.toml"
	rm -rf "src" "Cargo.lock"
fi

echo "Running benchmarks in: $TARGET_DIR with $RUNS warmup runs"

echo "🔵 - Go"
if [ ! -f "main.go" ]; then
	echo "❌ - I Gon't"
else
	#  Warmup runs
	./advent_code_2024_go
	{ time ./advent_code_2024_go; } 2>&1 | grep real | awk '{printf "🚀 Single - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do ./advent_code_2024_go >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total - %s\n", $2}'
fi

echo "🟠 - Rust"
if [ ! -f "main.rs" ]; then
	echo "❌ - I Rusn't"
else
	#  Warmup runs
	./target/release/advent_code_2024_rust
	{ time ./target/release/advent_code_2024_rust; } 2>&1 | grep real | awk '{printf "🚀 Single - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do ./target/release/advent_code_2024_rust >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total - %s\n", $2}'
fi

echo "⚪️ - Bun"
if [ ! -f "index.js" ]; then
	echo "❌ - I Bun't"
else
	#  Warmup runs
	bun run index.js
	{ time bun run index.js; } 2>&1 | grep real | awk '{printf "🚀 Single - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do bun run index.js >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total - %s\n", $2}'
fi

echo "🟢 - Node"
if [ ! -f "index.js" ]; then
	echo "❌ - I Noden't"
else
	#  Warmup runs
	node index.js
	{ time node index.js; } 2>&1 | grep real | awk '{printf "🚀 Single - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do node index.js >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total - %s\n", $2}'
fi
