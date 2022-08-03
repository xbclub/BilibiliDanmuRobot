package utiles

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	rainbowMode   bool
	trueColorMode bool
	spread        float64
	freq          float64
	message       string
)

const (
	blackBlock = "\033[40m  \033[0m"
	whiteBlock = "\033[47m  \033[0m"
)

//func parseFlags() {
//	flag.BoolVar(&rainbowMode, "r", false, "Rainbow mode")
//	flag.BoolVar(&trueColorMode, "t", false, "True color mode")
//	flag.Float64Var(&spread, "s", 3.0, "Rainbow `spread`")
//	flag.Float64Var(&freq, "f", 0.1, "Rainbow `frequency`")
//	flag.Usage = printUsage
//	flag.Parse()
//	message = flag.Arg(0)
//}

func detectTrueColorMode() bool {
	return os.Getenv("COLORTERM") == "truecolor"
}

// rainbow Reference: https://github.com/busyloop/lolcat/blob/b7ce4bd8882d22ee3db4b7d4d0df43eab6851cf5/lib/lolcat/lol.rb#L36
func rainbow(freq, i float64) (int, int, int) {
	red := int(math.Sin(freq*i+0)*128 + 128)
	green := int(math.Sin(freq*i+2*math.Pi/3)*127 + 128)
	blue := int(math.Sin(freq*i+4*math.Pi/3)*127 + 128)
	return red, green, blue
}

// rgbTo256 Reference: https://github.com/janlelis/paint/blob/7a76dc317a0d723f0acbf9dd393eb642822d6776/lib/paint.rb#L209
func rgbTo256(red, green, blue int, content string) string {
	var gray bool
	grayPossible := true
	sep := 42.5
	for grayPossible {
		if float64(red) < sep || float64(green) < sep || float64(blue) < sep {
			gray = float64(red) < sep && float64(green) < sep && float64(blue) < sep
			grayPossible = false
		}
		sep += 42.5
	}
	if gray {
		return fmt.Sprintf("\033[48;5;%dm%s\033[0m", 232+int(math.Round((float64(red)+float64(green)+float64(blue))/33)), content)
	}
	return fmt.Sprintf("\033[48;5;%dm%s\033[0m", 16+(int(6*float64(red)/256)*36+int(6*float64(green)/256)*6+int(6*float64(blue)/256)*1), content)
}

func rgbToTrueColor(red, green, blue int, content string) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", red, green, blue, content)
}

//func printUsage() {
//	fmt.Println("Usage: GoTerminalQRCode [OPTIONS]... <message>")
//	fmt.Println()
//	fmt.Println("Options:")
//	flag.PrintDefaults()
//	fmt.Println()
//	fmt.Println("Examples:")
//	fmt.Println("  $ GoTerminalQRCode Hello")
//	fmt.Println("  $ GoTerminalQRCode \"Hello World\"")
//	fmt.Println("  $ echo -n \"Hello World\" | GoTerminalQRCode")
//}

func GenerateQr(message string) error {
	//parseFlags()
	// When received message from pipeline
	if runtime.GOOS != "windows" && !terminal.IsTerminal(0) {
		pipeBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			//fmt.Fprintln(os.Stderr, "Failed to read Stdin:", err.Error())
			return err
		}
		message = string(pipeBytes)
	}
	if message == "" {
		//fmt.Fprintln(os.Stderr, "Message is empty")
		return errors.New("Message is empty")
	}
	//fmt.Println("Message:", message)
	if !trueColorMode {
		trueColorMode = detectTrueColorMode()
	}
	qr, err := qrcode.New(message, qrcode.Medium)
	if err != nil {
		//logx.Errorf(os.Stderr, "Failed to generate QR code:", err.Error())
		return err
	}
	var buf bytes.Buffer
	qrBitmap := qr.Bitmap()
	rainbowRandomSeed := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)
	rainbowOffset := float64(0)
	for y := range qrBitmap {
		for x := range qrBitmap[y] {
			if qrBitmap[y][x] {
				// Foreground
				if rainbowMode {
					red, green, blue := rainbow(freq, float64(rainbowRandomSeed)+(rainbowOffset/spread))
					if trueColorMode {
						buf.WriteString(rgbToTrueColor(red, green, blue, " "))
					} else {
						buf.WriteString(rgbTo256(red, green, blue, " "))
					}
					rainbowOffset++
					red, green, blue = rainbow(freq, float64(rainbowRandomSeed)+(rainbowOffset/spread))
					if trueColorMode {
						buf.WriteString(rgbToTrueColor(red, green, blue, " "))
					} else {
						buf.WriteString(rgbTo256(red, green, blue, " "))
					}
				} else {
					buf.WriteString(blackBlock)
				}
			} else {
				// Background
				if trueColorMode {
					buf.WriteString(rgbToTrueColor(255, 255, 255, "  "))
				} else {
					buf.WriteString(whiteBlock)
				}
				if rainbowMode {
					rainbowOffset++
				}
			}
			if rainbowMode {
				rainbowOffset++
			}
		}
		if rainbowMode {
			rainbowOffset = 0
			rainbowRandomSeed++
		}
		buf.WriteString("\n")
	}
	fmt.Println(buf.String())
	return nil
}
