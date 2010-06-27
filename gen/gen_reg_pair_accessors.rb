# AF BC DE HL AF' BC' DE' HL' IX IY SP PC

[['A', 'F'], ['B', 'C'], ['D', 'E'], ['A', 'F_'], ['B', 'C_'], ['D', 'E_'], ['H', 'L_']].each do |pair|
  prime = true if pair.last =~ /_/
  
  print <<CODE
func (z80 *Z80) #{pair}() uint16 {
  return uint16(z80.#{pair.last.downcase}) | (uint16(z80.#{prime ? pair.first.downcase + '_' : pair.first.downcase}) << 8)
}
CODE
end
