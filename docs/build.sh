#!/bin/sh

rm -rf wasm 
mkdir -p wasm

GOOS=js GOARCH=wasm go build -o wasm ../examples/...

cat > index.html << EOF
<!DOCTYPE html>
<html>
	<head>
		<title>gorge examples</title>
	</head>
	<body>
	Examples:
	<ul>
EOF
(cd wasm;
	for i in *; do
		echo "<li><a href=\"./wasm/?t=$i\">$i</a></li>" >> ../index.html
	done
)
cat >> index.html << EOF
	</ul>
	</body>
</html>
EOF

cp assets/wasm.html wasm/index.html





