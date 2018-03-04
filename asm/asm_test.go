package asm

import (
	"bytes"
	"testing"
)

func assemble(t *testing.T, code string) []byte {
	r := bytes.NewReader([]byte(code))
	result, err := Assemble(r, false)
	if err != nil {
		t.Error(err)
		return []byte{}
	}
	return result.Code
}

func fromHex(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return 0
	}
}

func checkASM(t *testing.T, asm string, expected string) {
	code := assemble(t, asm)

	b := make([]byte, len(code)*2)
	for i, j := 0, 0; i < len(code); i, j = i+1, j+2 {
		v := code[i]
		b[j+0] = hex[v>>4]
		b[j+1] = hex[v&0x0f]
	}
	s := string(b)

	if s != expected {
		t.Error("code doesn't match expected")
		t.Errorf("got: %s\n", s)
		t.Errorf("exp: %s\n", expected)
	}
}

func TestAddressingIMM(t *testing.T) {
	asm := `
	LDA #$20
	LDX #$20
	LDY #$20
	ADC #$20
	SBC #$20
	CMP #$20
	CPX #$20
	CPY #$20
	AND #$20
	ORA #$20
	EOR #$20`

	checkASM(t, asm, "A920A220A0206920E920C920E020C020292009204920")
}

func TestAddressingABS(t *testing.T) {
	asm := `
	LDA $2000
	LDX $2000
	LDY $2000
	STA $2000
	STX $2000
	STY $2000
	ADC $2000
	SBC $2000
	CMP $2000
	CPX $2000
	CPY $2000
	BIT $2000
	AND $2000
	ORA $2000
	EOR $2000
	INC $2000
	DEC $2000
	JMP $2000
	JSR $2000
	ASL $2000
	LSR $2000
	ROL $2000
	ROR $2000
	LDA A:$20
	LDA ABS:$20`

	checkASM(t, asm, "AD0020AE0020AC00208D00208E00208C00206D0020ED0020CD0020"+
		"EC0020CC00202C00202D00200D00204D0020EE0020CE00204C00202000200E0020"+
		"4E00202E00206E0020AD2000AD2000")
}

func TestAddressingABX(t *testing.T) {
	asm := `
	LDA $2000,X
	LDY $2000,X
	STA $2000,X
	ADC $2000,X
	SBC $2000,X
	CMP $2000,X
	AND $2000,X
	ORA $2000,X
	EOR $2000,X
	INC $2000,X
	DEC $2000,X
	ASL $2000,X
	LSR $2000,X
	ROL $2000,X
	ROR $2000,X`

	checkASM(t, asm, "BD0020BC00209D00207D0020FD0020DD00203D00201D00205D0020"+
		"FE0020DE00201E00205E00203E00207E0020")
}

func TestAddressingABY(t *testing.T) {
	asm := `
	LDA $2000,Y
	LDX $2000,Y
	STA $2000,Y
	ADC $2000,Y
	SBC $2000,Y
	CMP $2000,Y
	AND $2000,Y
	ORA $2000,Y
	EOR $2000,Y`

	checkASM(t, asm, "B90020BE0020990020790020F90020D90020390020190020590020")
}

func TestAddressingZPG(t *testing.T) {
	asm := `
	LDA $20
	LDX $20
	LDY $20
	STA $20
	STX $20
	STY $20
	ADC $20
	SBC $20
	CMP $20
	CPX $20
	CPY $20
	BIT $20
	AND $20
	ORA $20
	EOR $20
	INC $20
	DEC $20
	ASL $20
	LSR $20
	ROL $20
	ROR $20`

	checkASM(t, asm, "A520A620A4208520862084206520E520C520E420C42024202520"+
		"05204520E620C6200620462026206620")
}

func TestAddressingIND(t *testing.T) {
	asm := `
	JMP ($20)
	JMP ($2000)`

	checkASM(t, asm, "6C20006C0020")
}