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

	# Build comparison Go files
	for file in alt-*.go; do
		if [ -f "$file" ]; then
			output_name="advent_code_2024_go_${file%.*}"
			go build -o "$output_name" "$file"
		fi
	done
fi

# Rust build
if [ -f "main.rs" ]; then
	echo "Building Rust executable..."
	mkdir -p "src"
	cp "main.rs" "src/main.rs"

	CARGO_TOML_TEMPLATE='[package]
name = "_NAME_"
version = "0.1.0"
edition = "2021"

[dependencies]
tokio = { version = "1.28", features = ["full"] }
regex = "1.7"
itertools = "0.10.5"'

	# Copy comparison Rust files
	for file in alt-*.rs; do
		if [ -f "$file" ]; then
			cp "$file" "src/main.rs"
			# Generate Cargo.toml with templated content
			echo "${CARGO_TOML_TEMPLATE//_NAME_/advent_code_2024_rust_${file%.*}}" >"Cargo.toml"
			cargo build --release
			mv "./target/release/advent_code_2024_rust_${file%.*}" .
		fi
	done

	# Build main Rust file
	cp "main.rs" "src/main.rs"
	echo "${CARGO_TOML_TEMPLATE//_NAME_/advent_code_2024_rust}" >"Cargo.toml"
	cargo build --release
	rm -rf "src" "Cargo.toml" "Cargo.lock"
fi

echo "------------------------------------------------"
echo "> Running benchmarks for [Day $TARGET_DIR] with [$RUNS runs]"
echo "------------------------------------------------"

# Go benchmarks
echo "🔵 - Go"
if [ ! -f "main.go" ]; then
	echo "❌ - I Gon't"
else
	echo "[Default]"
	./advent_code_2024_go # Warmup
	{ time ./advent_code_2024_go; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do ./advent_code_2024_go >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'

	# Run comparison Go files
	for file in advent_code_2024_go_alt-*; do
		if [ -f "$file" ]; then
			echo "[Alternative (${file#advent_code_2024_go_})]"
			./"$file" # Warmup
			{ time ./"$file"; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
			{ time for i in $(seq 1 $RUNS); do ./"$file" >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'
		fi
	done
fi

echo "------------------------------"
echo "🟠 - Rust"
if [ ! -f "main.rs" ]; then
	echo "❌ - I Rusn't"
else
	echo "[Default]"
	./target/release/advent_code_2024_rust # Warmup
	{ time ./target/release/advent_code_2024_rust; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do ./target/release/advent_code_2024_rust >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'

	# Run comparison Rust files
	for file in advent_code_2024_rust_alt-*; do
		if [ -f "$file" ]; then
			echo "[Alternative (${file#advent_code_2024_rust_})]"
			./"$file" # Warmup
			{ time ./"$file"; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
			{ time for i in $(seq 1 $RUNS); do ./"$file" >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'
		fi
	done
fi

echo "------------------------------"
echo "⚪️ - Bun"
if [ ! -f "index.js" ]; then
	echo "❌ - I Bun't"
else
	echo "[Default]"
	bun run index.js # Warmup
	{ time bun run index.js; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do bun run index.js >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'

	# Run comparison JS files with Bun
	for file in alt-*.js; do
		if [ -f "$file" ]; then
			echo "[Alternative (${file})]"
			bun run "$file" # Warmup
			{ time bun run "$file"; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
			{ time for i in $(seq 1 $RUNS); do bun run "$file" >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'
		fi
	done
fi

echo "------------------------------"
echo "🟢 - Node"
if [ ! -f "index.js" ]; then
	echo "❌ - I Noden't"
else
	echo "[Default]"
	node index.js # Warmup
	{ time node index.js; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
	{ time for i in $(seq 1 $RUNS); do node index.js >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'

	# Run comparison JS files with Node
	for file in alt-*.js; do
		if [ -f "$file" ]; then
			echo "[Alternative (${file})]"
			node "$file" # Warmup
			{ time node "$file"; } 2>&1 | grep real | awk '{printf "🚀 Single run - %s\n", $2}'
			{ time for i in $(seq 1 $RUNS); do node "$file" >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{printf "🚀 Total all runs - %s\n", $2}'
		fi
	done
fi
echo "------------------------------"
