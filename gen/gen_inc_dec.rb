%w{ a b c d e h l }.each do |reg|
  print <<CODE
func (z80 *Z80) inc#{reg.upcase}() {
	z80.#{reg}++
	z80.f = ( z80.f & FLAG_C ) | ( ternOpB(z80.#{reg} == 0x80, FLAG_V, 0) ) | ( ternOpB((z80.#{reg} & 0x0f) == 1, 0, FLAG_H) ) | z80.sz53Table[z80.#{reg}]
}

func (z80 *Z80) dec#{reg.upcase}() {
	z80.f = ( z80.f & FLAG_C ) | ( ternOpB(z80.#{reg} & 0x0f == 1, 0, FLAG_H )) | FLAG_N
	z80.#{reg}--
	z80.f |= ( ternOpB(z80.#{reg} == 0x7f, FLAG_V, 0) ) | z80.sz53Table[z80.#{reg}]

}

CODE
end
