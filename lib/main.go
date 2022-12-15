//go:build js && wasm

package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"syscall/js"
)

const (
	colorMax = 255

	colorBlack = "#000"
	colorWhite = "#fff"
)

// マンデルブロ集合の漸化式を再帰的に計算する回数
var recursionLimit int

func main() {
	js.Global().Set("calculate", js.FuncOf(calculate))
	select {}
}

func calculate(this js.Value, args []js.Value) any {
	// 引数が期待した個数無かったらすぐ終わり
	if len(args) != 6 {
		return js.ValueOf([]any{})
	}

	// 適当に引数をパース
	recursionLimit = args[0].Int()
	scale := args[1].Float()
	width := args[2].Int()
	height := args[3].Int()
	offsetWidth := args[4].Float()
	offsetHeight := args[5].Float()

	// result[x][y] = 結果 というような行列になる
	result := make([]any, 0, width)
	for x := 0; x < width; x++ {
		row := make([]any, 0, height)
		for y := 0; y < height; y++ {
			a := float64(x) / float64(width/2)
			b := float64(y) / float64(height/2)

			// 描画するキャンバスは(0,0)を左上とする座標になる
			// 描画しない範囲にあることが事前にわかるのであれば計算をスキップする
			if a < 0 && y > 0 {
				row = append(row, []any{colorWhite, false})
				continue
			}

			// 表示の都合で良い感じに座標じオフセットをかける
			// 拡大、縮小もする
			a = (a - offsetWidth) / scale
			b = (b - offsetHeight) / scale

			depth, ok := mandelbro(0, 0, complex(a, b))
			// 発散するまでに掛かった回数が多いほど色が濃くなる
			row = append(row, []any{colorize(depth), ok})
		}
		result = append(result, row)
	}

	return js.ValueOf(result)
}

// c = a + bi
// z(0) = 0
// z(n+1) = z(n)^2 + c
// マンデルブロ集合の漸化式を再帰的に解く
// 規定の回数でも発散しなかった場合はマンデルブロ集合に含まれる
// 含まれる場合は第二返り値がtrueになる
func mandelbro(depth int, z complex128, c complex128) (_ int, ok bool) {
	z = z*z + c
	if cmplx.Abs(z) > 2 {
		// zの絶対値が2を越えるなら発散するとみなす
		return depth, false
	}
	if depth < recursionLimit {
		// 規定の回数を超えてなかったらn+1を計算する
		return mandelbro(depth+1, z, c)
	}
	return depth, true
}

func colorize(depth int) string {
	// 規定の回数まで到達していたら真っ黒
	if depth == recursionLimit {
		return colorBlack
	}
	// 1回目で発散したら真っ白
	if depth == 0 {
		return colorWhite
	}
	// 適当に発散までに掛かった回数を濃度(数学用語ではない)として良い感じに計算する
	// かかればかかるほど色が黒くなるようにする
	cardinality := math.Pow(float64(depth)/float64(recursionLimit), math.Sqrt2)
	color := colorMax - int(math.Floor(colorMax*cardinality))
	return fmt.Sprintf("#%x%x%x", color, color, color)
}
