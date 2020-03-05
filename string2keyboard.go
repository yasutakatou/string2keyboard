package string2keyboard

import (
	"runtime"
	"time"
	"strings"
	"strconv"

	"github.com/micmonay/keybd_event"
)

type keySet struct {
	code  int
	shift bool
}

//KeyboardWrite emulate keyboard input from string
func KeyboardWrite(textInput string, fCtrl bool, fAlt bool,LiveRawcodeChar string) error {
	byteRawcode := []byte(LiveRawcodeChar)[0]

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	// For linux, it is very important wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	if fCtrl == true {
		kb.HasCTRL(true)
	} else {
		kb.HasCTRL(false)
	}

	if fAlt == true {
		kb.HasALT(true)
	} else {
		kb.HasALT(false)
	}

	//Should we skip next character in string
	//Used if we found some escape sequence
	skip := false
	for i := 0; i < len(textInput); i++ {
		c := textInput[i]

		if !skip {
			if c != '\\' {
				kb.SetKeys(names[string(c)].code)
				kb.HasSHIFT(names[string(c)].shift)
			} else {
				//Found backslash escape character
				//Check next character
				switch textInput[i+1] {
				case 'n':
					//Found newline character sequence
					kb.SetKeys(names["ENTER"].code)
					skip = true
				case '\\':
					//Found backslash character sequence
					kb.SetKeys(names["\\"].code)
					kb.HasSHIFT(names["\\"].shift)
					skip = true
				case 'b':
					//Found backspace character sequence
					kb.SetKeys(names["BACKSPACE"].code)
					skip = true
				case 't':
					//Found tab character sequence
					kb.SetKeys(names["TAB"].code)
					skip = true
				case '"':
					//Found double quote character sequence
					kb.SetKeys(names["\""].code)
					kb.HasSHIFT(names["\""].shift)
					skip = true
				case byteRawcode:
					rawValue := textInput[i+2:i+2+strings.Index(textInput[i+2:], LiveRawcodeChar)]
					cnt, err := strconv.Atoi(asciiToGo(rawValue))
					if err != nil { return err }
					kb.SetKeys(cnt)

					i += 2 + strings.Index(textInput[i+2:], LiveRawcodeChar)
					skip = true
				default:
					//Nothing special, jsut backslash output
					kb.SetKeys(names["\\"].code)
					kb.HasSHIFT(names["\\"].shift)
				}

			}
			err = kb.Launching()
			if err != nil {
				return err
			}
		} else {
			skip = false
		}

	}
	return nil

}

/*
	F1		112	VK_F1			59
	F2		113	VK_F2			60
	F3		114	VK_F3			61
	F4		115	VK_F4			62
	F5		116	VK_F5			63
	F6		117	VK_F6			64
	F7		118	VK_F7			65
	F8		119	VK_F8			66
	F9		120	VK_F9			67
	F10		121	VK_F10			68
	F11		122	VK_F11			87
	F12		123	VK_F12			88
	↑		38	VK_UP			4133
	↓		40	VK_DOWN			4135
	←		37	VK_LEFT 		4132
	→		39	VK_RIGHT	 	4134
	Esc		27	VK_ESC			1
	Capa	20	VK_CAPSLOCK		58
	NumLock	144	VK_NUMLOCK		69
	Insert	45	VK_INSERT		4140
	Delete	46	VK_DELETE 		4140
	Home	36	VK_HOME  		4131
	End		35	VK_END     		4130
	PageUp	33	VK_PAGEUP 		4128
	PageDw	34	VK_PAGEDOWN 	4128
	SLock	145	VK_SCROLLLOCK	70
*/

var asciiToGoMap = map[string]string {
	"112": "59",
	"113": "60",
	"114": "61",
	"115": "62",
	"116": "63",
	"117": "64",
	"118": "65",
	"119": "66",
	"120": "67",
	"121": "68",
	"122": "87",
	"123": "88",
	"38": "4133",
	"40": "4135",
	"37": "4132",
	"39": "4134",
	"27": "1",
	"20": "58",
	"144": "69",
	"45": "4140",
	"46": "4140",
	"36": "4131",
	"35": "4130",
	"33": "4128",
	"34": "4128",
	"145": "70",	
}

func asciiToGo(ascii string) string {
    return asciiToGoMap[ascii]
}
