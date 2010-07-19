#!/usr/bin/perl -w

# z80.pl: generate C code for Z80 opcodes
# $Id: z80.pl 3681 2008-06-16 09:40:29Z pak21 $

# Copyright (c) 1999-2006 Philip Kendall

# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License along
# with this program; if not, write to the Free Software Foundation, Inc.,
# 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

# Author contact information:

# E-mail: philip-fuse@shadowmagic.org.uk

use strict;

use Fuse;

# The status of which flags relates to which condition

# These conditions involve !( F & FLAG_<whatever> )
my %not = map { $_ => 1 } qw( NC NZ P PO );

# Use F & FLAG_<whatever>
my %flag = (

      C => 'C', NC => 'C',
     PE => 'P', PO => 'P',
      M => 'S',  P => 'S',
      Z => 'Z', NZ => 'Z',

);

# Generalised opcode routines

sub arithmetic_logical ($$$) {

    my( $opcode, $arg1, $arg2 ) = @_;

    unless( $arg2 ) { $arg2 = $arg1; $arg1 = 'A'; }

    if( length $arg1 == 1 ) {
	if( length $arg2 == 1 or $arg2 =~ /^REGISTER[HL]$/ ) {
            my($lcopcode, $lcarg2); $lcopcode = lc($opcode); $lcarg2 = lc($arg2);
	    print "      z80.$lcopcode(z80.$lcarg2)\n";
	} elsif( $arg2 eq '(REGISTER+dd)' ) {
            my($lcopcode, $lcarg2); $lcopcode = lc($opcode);
	    print << "CODE";
      
	var offset, bytetemp byte
	offset = z80.memory.readByte( z80.pc )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++
	bytetemp = z80.memory.readByte(uint16(int(z80.REGISTER()) + int(signExtend(offset))))
	z80.$lcopcode(bytetemp)

CODE
	} else {
	    my $register = ( $arg2 eq '(HL)' ? 'HL' : 'PC' );
	    my $increment = ( $register eq 'PC' ? 'z80.pc++' : '' );
	    my $lcopcode = lc($opcode);
	    print << "CODE";
      {
	var bytetemp byte = z80.memory.readByte(z80.$register())
        $increment
	z80.$lcopcode(bytetemp)
      }
CODE
	}
    } elsif( $opcode eq 'ADD' ) {
	$opcode = lc($opcode);
	$arg1 = lc($arg1);
	print << "CODE";
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.${opcode}16(z80.$arg1, z80.$arg2());
CODE
    } elsif( $arg1 eq 'HL' and length $arg2 == 2 ) {
	$opcode = lc($opcode);
	print << "CODE";
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.${opcode}16(z80.$arg2());
CODE
    }
}

sub call_jp ($$$) {

    my( $opcode, $condition, $offset ) = @_;
    my $lcopcode = lc($opcode);
    if( not defined $offset ) {
	print "      z80.$lcopcode()\n";
    } else {
	my $condition_string;
	if( defined $not{$condition} ) {
	    $condition_string = "(z80.f & FLAG_$flag{$condition}) == 0";
	} else {
	    $condition_string = "(z80.f & FLAG_$flag{$condition}) != 0";
	}
	print << "CALL";
      if($condition_string) {
	z80.$lcopcode()
      } else {
	z80.memory.contendRead(z80.pc, 3); z80.memory.contendRead( z80.pc + 1, 3 ); z80.pc += 2;
      }
CALL
    }
}

sub cpi_cpd ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'CPI' ? 'inc' : 'dec' );

    print << "CODE";
      
	var value, bytetemp, lookup byte 

	value = z80.memory.readByte( z80.HL() )
	bytetemp = z80.a - value
        lookup = ((z80.a & 0x08 ) >> 3 ) | (((value) & 0x08 ) >> 2 ) | ((bytetemp & 0x08 ) >> 1)

	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.${modifier}HL(); z80.decBC();
	z80.f = (z80.f & FLAG_C) | ternOpB(z80.BC() != 0, FLAG_V | FLAG_N, FLAG_N) | halfcarrySubTable[lookup] | ternOpB(bytetemp != 0, 0, FLAG_Z) | (bytetemp & FLAG_S )
	if((z80.f & FLAG_H) != 0) { bytetemp-- }
	z80.f |= (bytetemp & FLAG_3) | ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)
      
CODE
}

sub cpir_cpdr ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'CPIR' ? 'inc' : 'dec' );

    print << "CODE";
	var value, bytetemp, lookup byte

	value = z80.memory.readByte( z80.HL() )
	bytetemp = z80.a - value
        lookup = ((z80.a & 0x08) >> 3) | (((value) & 0x08) >> 2) | ((bytetemp & 0x08) >> 1)

	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.decBC()
	z80.f = ( z80.f & FLAG_C ) | ( ternOpB(z80.BC() != 0, ( FLAG_V | FLAG_N ),FLAG_N)) | halfcarrySubTable[lookup] | ( ternOpB(bytetemp != 0, 0, FLAG_Z )) | ( bytetemp & FLAG_S );
	if((z80.f & FLAG_H) != 0) {
	    bytetemp--
        }
	z80.f |= ( bytetemp & FLAG_3 ) | ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)
	if( ( z80.f & ( FLAG_V | FLAG_Z ) ) == FLAG_V ) {
	  z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 );
	  z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 );
	  z80.memory.contendReadNoMreq( z80.HL(), 1 );
	  z80.pc-=2;
	}
	z80.${modifier}HL()
CODE
}

sub inc_dec ($$) {

    my( $opcode, $arg ) = @_;

    my $modifier = ( $opcode eq 'INC' ? 'inc' : 'dec' );

    if( length $arg == 1 or $arg =~ /^REGISTER[HL]$/ ) {
	print "      z80.".lc($opcode)."$arg()\n";
    } elsif( length $arg == 2 or $arg eq 'REGISTER' ) {
	print << "CODE";
	z80.memory.contendReadNoMreq( z80.IR(), 1 )
	z80.memory.contendReadNoMreq( z80.IR(), 1 )
	z80.$modifier${arg}()
CODE
    } elsif( $arg eq '(HL)' ) {
	$opcode = lc($opcode);
	print << "CODE";
      {
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.$opcode(&bytetemp)
	z80.memory.writeByte(z80.HL(), bytetemp)
      }
CODE
    } elsif( $arg eq '(REGISTER+dd)' ) {
	my $lcopcode = lc($opcode);
	print << "CODE";
	var offset, bytetemp byte
	var wordtemp uint16
	offset = z80.memory.readByte( z80.pc )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++
	wordtemp = uint16(int(z80.REGISTER()) + int(signExtend(offset)))
	bytetemp = z80.memory.readByte( wordtemp )
	z80.memory.contendReadNoMreq( wordtemp, 1 )
	z80.$lcopcode(&bytetemp)
	z80.memory.writeByte(wordtemp,bytetemp)
CODE
    }

}

sub ini_ind ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'INI' ? 'inc' : 'dec' );
    my $operation = ( $opcode eq 'INI' ? '+' : '-' );
    print << "CODE";
	var initemp, initemp2 byte

	z80.memory.contendReadNoMreq( z80.IR(), 1 );
	initemp = z80.readPort(z80.BC());
	z80.memory.writeByte( z80.HL(), initemp );

        z80.b--; z80.${modifier}HL()
        initemp2 = initemp + z80.c $operation 1;
	z80.f = ternOpB((initemp & 0x80) != 0, FLAG_N, 0) | ternOpB(initemp2 < initemp, FLAG_H | FLAG_C, 0) | ternOpB(z80.parityTable[(initemp2 & 0x07) ^ z80.b] != 0, FLAG_P, 0 ) | z80.sz53Table[z80.b]
CODE
}

sub inir_indr ($) {

    my( $opcode ) = @_;

    my $operation = ( $opcode eq 'INIR' ? '+' : '-' );
    my $modifier = ( $opcode eq 'INIR' ? 'inc' : 'dec' );

    print << "CODE";
	var initemp, initemp2 byte

	z80.memory.contendReadNoMreq( z80.IR(), 1 );
	initemp = z80.readPort(z80.BC());
	z80.memory.writeByte( z80.HL(), initemp );

	z80.b--;
        initemp2 = initemp + z80.c $operation 1;
	z80.f = ternOpB(initemp & 0x80 != 0, FLAG_N, 0) |
                ternOpB(initemp2 < initemp, FLAG_H | FLAG_C, 0 ) |
                ternOpB(z80.parityTable[ ( initemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0) |
                z80.sz53Table[z80.b];

	if( z80.b != 0 ) {
	  z80.memory.contendWriteNoMreq( z80.HL(), 1 ); z80.memory.contendWriteNoMreq( z80.HL(), 1 );
	  z80.memory.contendWriteNoMreq( z80.HL(), 1 ); z80.memory.contendWriteNoMreq( z80.HL(), 1 );
	  z80.memory.contendWriteNoMreq( z80.HL(), 1 );
	  z80.pc -= 2;
	}
        z80.${modifier}HL()
CODE
}


sub ldi_ldd ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'LDI' ? 'inc' : 'dec' );

    print << "CODE";
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.decBC()
	z80.memory.writeByte(z80.DE(), bytetemp);
	z80.memory.contendWriteNoMreq( z80.DE(), 1 ); z80.memory.contendWriteNoMreq( z80.DE(), 1 );
	z80.${modifier}DE(); z80.${modifier}HL();
	bytetemp += z80.a;
	z80.f = ( z80.f & ( FLAG_C | FLAG_Z | FLAG_S ) ) | ternOpB(z80.BC() != 0, FLAG_V, 0) |
	  ( bytetemp & FLAG_3 ) | ternOpB((bytetemp & 0x02) != 0, FLAG_5, 0)
CODE
}

sub ldir_lddr ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'LDIR' ? 'inc' : 'dec' );

    print << "CODE";
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.memory.writeByte(z80.DE(), bytetemp);
	z80.memory.contendWriteNoMreq(z80.DE(), 1); z80.memory.contendWriteNoMreq(z80.DE(), 1 );
	z80.decBC()
	bytetemp += z80.a;
	z80.f = (z80.f & ( FLAG_C | FLAG_Z | FLAG_S )) | ternOpB(z80.BC() != 0, FLAG_V, 0 ) | (bytetemp & FLAG_3) | ternOpB((bytetemp & 0x02 != 0), FLAG_5, 0 )
	if(z80.BC() != 0) {
	  z80.memory.contendWriteNoMreq( z80.DE(), 1 ); z80.memory.contendWriteNoMreq( z80.DE(), 1 );
	  z80.memory.contendWriteNoMreq( z80.DE(), 1 ); z80.memory.contendWriteNoMreq( z80.DE(), 1 );
	  z80.memory.contendWriteNoMreq( z80.DE(), 1 );
	  z80.pc -= 2
	}
        z80.${modifier}HL(); z80.${modifier}DE()
CODE
}

sub otir_otdr ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'OTIR' ? 'inc' : 'dec' );

    print << "CODE";
	var outitemp, outitemp2 byte

	z80.memory.contendReadNoMreq( z80.IR(), 1 );
	outitemp = z80.memory.readByte( z80.HL() );
	z80.b--;	/* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp);

	z80.${modifier}HL()
        outitemp2 = outitemp + z80.l;
	z80.f = ternOpB((outitemp & 0x80) != 0, FLAG_N, 0 ) |
            ternOpB(outitemp2 < outitemp, FLAG_H | FLAG_C, 0) |
            ternOpB(z80.parityTable[ ( outitemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0 ) |
            z80.sz53Table[z80.b]

	if( z80.b != 0 ) {
	  z80.memory.contendReadNoMreq( z80.BC(), 1 ); z80.memory.contendReadNoMreq( z80.BC(), 1 );
	  z80.memory.contendReadNoMreq( z80.BC(), 1 ); z80.memory.contendReadNoMreq( z80.BC(), 1 );
	  z80.memory.contendReadNoMreq( z80.BC(), 1 );
	  z80.pc -= 2;
	}
CODE
}

sub outi_outd ($) {

    my( $opcode ) = @_;

    my $modifier = ( $opcode eq 'OUTI' ? 'inc' : 'dec' );

    print << "CODE";
	var outitemp, outitemp2 byte

	z80.memory.contendReadNoMreq( z80.IR(), 1 );
	outitemp = z80.memory.readByte( z80.HL() );
	z80.b--;	/* This does happen first, despite what the specs say */
	z80.writePort(z80.BC(), outitemp);

	z80.${modifier}HL()
        outitemp2 = outitemp + z80.l;
	z80.f = ternOpB((outitemp & 0x80) != 0, FLAG_N, 0) |
            ternOpB(outitemp2 < outitemp, FLAG_H | FLAG_C, 0) |
            ternOpB(z80.parityTable[ ( outitemp2 & 0x07 ) ^ z80.b ] != 0, FLAG_P, 0 ) |
            z80.sz53Table[z80.b];
CODE
}

sub push_pop ($$) {

    my( $opcode, $regpair ) = @_;

    my( $high, $low );

    if( $regpair eq 'REGISTER' ) {
	( $high, $low ) = ( 'REGISTERH', 'REGISTERL' );
    } else {
	( $high, $low ) = ( $regpair =~ /^(.)(.)$/ );
    }
    my $lcopcode = lc($opcode);
    my $lclow = lc($low);
    my $lchigh = lc($high);
    if( $lcopcode eq 'pop') {
	print "      z80.${lcopcode}16(&z80.$lclow, &z80.$lchigh)\n";
    } else {
	print "      z80.${lcopcode}16(z80.$lclow, z80.$lchigh)\n";
    }
}

sub res_set_hexmask ($$) {

    my( $opcode, $bit ) = @_;

    my $mask = 1 << $bit;
    $mask = 0xff - $mask if $opcode eq 'RES';

    sprintf '0x%02x', $mask;
}

sub res_set ($$$) {

    my( $opcode, $bit, $register ) = @_;

    my $operator = ( $opcode eq 'RES' ? '&' : '|' );

    my $hex_mask = res_set_hexmask( $opcode, $bit );

    my $lcopcode = lc($opcode);
    my $lcregister = lc($register);

    if( length $register == 1 ) {
	print "      z80.$lcregister $operator= $hex_mask;\n";
    } elsif( $register eq '(HL)' ) {
	print << "CODE";
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.memory.contendReadNoMreq( z80.HL(), 1 );
	z80.memory.writeByte( z80.HL(), bytetemp $operator $hex_mask );
CODE
    } elsif( $register eq '(REGISTER+dd)' ) {
	print << "CODE";
   
	var bytetemp byte
	bytetemp = z80.memory.readByte( tempaddr );
	z80.memory.contendReadNoMreq( tempaddr, 1 );
	z80.memory.writeByte( tempaddr, bytetemp $operator $hex_mask );

CODE
    }
}

sub rotate_shift ($$) {

    my( $opcode, $register ) = @_;
    my $lcopcode = lc($opcode);
    my $lcregister = lc($register);

    if( length $register == 1 ) {
	print "      z80.$lcopcode(&z80.$lcregister)\n";
    } elsif( $register eq '(HL)' ) {
	print << "CODE";
	var bytetemp byte = z80.memory.readByte(z80.HL())
	z80.memory.contendReadNoMreq( z80.HL(), 1 );
	z80.$lcopcode(&bytetemp);
	z80.memory.writeByte(z80.HL(),bytetemp);
CODE
    } elsif( $register eq '(REGISTER+dd)' ) {
	print << "CODE";
	var bytetemp byte = z80.memory.readByte(tempaddr);
	z80.memory.contendReadNoMreq( tempaddr, 1 );
	z80.$lcopcode(&bytetemp);
	z80.memory.writeByte(tempaddr,bytetemp);
CODE
    }
}

# Individual opcode routines

sub opcode_ADC (@) { arithmetic_logical( 'ADC', $_[0], $_[1] ); }

sub opcode_ADD (@) { arithmetic_logical( 'ADD', $_[0], $_[1] ); }

sub opcode_AND (@) { arithmetic_logical( 'AND', $_[0], $_[1] ); }

sub opcode_BIT (@) {

    my( $bit, $register ) = @_;
    my $lcregister = lc($register);
    if( length $register == 1 ) {
	print "      z80.bit($bit, z80.$lcregister );\n";
    } elsif( $register eq '(REGISTER+dd)' ) {
	print << "BIT";
	bytetemp := z80.memory.readByte( tempaddr )
	z80.memory.contendReadNoMreq( tempaddr, 1 )
	z80.biti($bit, bytetemp, tempaddr)
BIT
    } else {
	print << "BIT";
	bytetemp := z80.memory.readByte( z80.HL() );
	z80.memory.contendReadNoMreq( z80.HL(), 1 );
	z80.bit($bit, bytetemp)
BIT
    }
}

sub opcode_CALL (@) { call_jp( 'CALL', $_[0], $_[1] ); }

sub opcode_CCF (@) {
    print << "CCF";
      z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |
	ternOpB( ( z80.f & FLAG_C ) != 0, FLAG_H, FLAG_C ) | ( z80.a & ( FLAG_3 | FLAG_5 ) );
CCF
}

sub opcode_CP (@) { arithmetic_logical( 'CP', $_[0], $_[1] ); }

sub opcode_CPD (@) { cpi_cpd( 'CPD' ); }

sub opcode_CPDR (@) { cpir_cpdr( 'CPDR' ); }

sub opcode_CPI (@) { cpi_cpd( 'CPI' ); }

sub opcode_CPIR (@) { cpir_cpdr( 'CPIR' ); }

sub opcode_CPL (@) {
    print << "CPL";
      z80.a ^= 0xff;
      z80.f = ( z80.f & ( FLAG_C | FLAG_P | FLAG_Z | FLAG_S ) ) |
	( z80.a & ( FLAG_3 | FLAG_5 ) ) | ( FLAG_N | FLAG_H );
CPL
}

sub opcode_DAA (@) {
    print << "DAA";
	var add, carry byte = 0, ( z80.f & FLAG_C )
        if( ( (z80.f & FLAG_H ) != 0) || ( ( z80.a & 0x0f ) > 9 ) ) { add = 6 }
        if( (carry != 0) || ( z80.a > 0x99 ) ) { add |= 0x60 }
        if( z80.a > 0x99 ) { carry = FLAG_C }
	if( (z80.f & FLAG_N) != 0 ) {
	  z80.sub(add)
	} else {
	  z80.add(add)
	}
        var temp int = (int(z80.f) & ^(FLAG_C | FLAG_P)) | int(carry) | int(z80.parityTable[z80.a])
	z80.f = byte(temp)
DAA
}

sub opcode_DEC (@) { inc_dec( 'DEC', $_[0] ); }

sub opcode_DI (@) { print "      z80.iff1, z80.iff2 = 0, 0\n"; }

sub opcode_DJNZ (@) {
    print << "DJNZ";
      z80.memory.contendReadNoMreq(z80.IR(), 1)
      z80.b--
      if(z80.b != 0) {
	z80.jr()
      } else {
	z80.memory.contendRead( z80.pc, 3 )
      }
      z80.pc++
DJNZ
}

sub opcode_EI (@) {
    print << "EI";
      /* Interrupts are not accepted immediately after an EI, but are
	 accepted after the next instruction */
      z80.iff1, z80.iff2 = 1, 1
      z80.interruptsEnabledAt = int(tstates)
      // eventAdd(tstates + 1, z80InterruptEvent)
EI
}

sub opcode_EX (@) {

    my( $arg1, $arg2 ) = @_;

    if( $arg1 eq 'AF' and $arg2 eq "AF'" ) {
	print << "EX";
      /* Tape saving trap: note this traps the EX AF,AF\' at #04d0, not
	 #04d1 as PC has already been incremented */
      /* 0x76 - Timex 2068 save routine in EXROM */
      if( z80.pc == 0x04d1 || z80.pc == 0x0077 ) {
	  if( z80.tapeSaveTrap() == 0 ) { break }
      }

      var olda, oldf = z80.a, z80.f
      z80.a = z80.a_; z80.f = z80.f_
      z80.a_ = olda; z80.f_ = oldf
EX
    } elsif( $arg1 eq '(SP)' and ( $arg2 eq 'HL' or $arg2 eq 'REGISTER' ) ) {

	my( $high, $low );

	if( $arg2 eq 'HL' ) {
	    ( $high, $low ) = qw( H L );
	} else {
	    ( $high, $low ) = qw( REGISTERH REGISTERL );
	}
	my $lclow = lc($low);
	my $lchigh = lc($high);

	print << "EX";
	var bytetempl, bytetemph byte
	bytetempl = z80.memory.readByte( z80.SP() )
	bytetemph = z80.memory.readByte( z80.SP() + 1 )
        z80.memory.contendReadNoMreq( z80.SP() + 1, 1 )
	z80.memory.writeByte( z80.SP() + 1, z80.$lchigh )
	z80.memory.writeByte( z80.SP(),     z80.$lclow  )
	z80.memory.contendWriteNoMreq( z80.SP(), 1 )
        z80.memory.contendWriteNoMreq( z80.SP(), 1 )
	z80.$lclow = bytetempl
        z80.$lchigh = bytetemph
EX
    } elsif( $arg1 eq 'DE' and $arg2 eq 'HL' ) {
	print << "EX";
	var wordtemp uint16 = z80.DE()
        z80.setDE(z80.HL())
        z80.setHL(wordtemp)
EX
    }
}

sub opcode_EXX (@) {
    print << "EXX";
	var wordtemp uint16
	wordtemp = z80.BC() 
        z80.setBC(z80.BC_())
        z80.setBC_(wordtemp)

	wordtemp = z80.DE() 
        z80.setDE(z80.DE_())
        z80.setDE_(wordtemp)

	wordtemp = z80.HL()
        z80.setHL(z80.HL_())
        z80.setHL_(wordtemp)
EXX
}

sub opcode_HALT (@) { print "      z80.halted=1;\n      z80.pc--;\n"; }

sub opcode_IM (@) {

    my( $mode ) = @_;

    print "      z80.im = $mode;\n";
}

sub opcode_IN (@) {

    my( $register, $port ) = @_;

    if( $register eq 'A' and $port eq '(nn)' ) {
	print << "IN";
	var intemp uint16
	intemp = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8 )
	z80.pc++
        z80.a = z80.readPort(intemp)
IN
    } elsif( $register eq 'F' and $port eq '(C)' ) {
	print << "IN";
	var bytetemp byte
	z80.in(&bytetemp, z80.BC())
IN
    } elsif( length $register == 1 and $port eq '(C)' ) {
	my $lcregister = lc($register);
	print << "IN";
	z80.in(&z80.$lcregister, z80.BC())
IN
    }
}

sub opcode_INC (@) { inc_dec( 'INC', $_[0] ); }

sub opcode_IND (@) { ini_ind( 'IND' ); }

sub opcode_INDR (@) { inir_indr( 'INDR' ); }

sub opcode_INI (@) { ini_ind( 'INI' ); }

sub opcode_INIR (@) { inir_indr( 'INIR' ); }

sub opcode_JP (@) {

    my( $condition, $offset ) = @_;

    if( $condition eq 'HL' or $condition eq 'REGISTER' ) {
	print "      z80.pc = z80.$condition()\t\t/* NB: NOT INDIRECT! */\n";
	return;
    } else {
	call_jp( 'JP', $condition, $offset );
    }
}

sub opcode_JR (@) {

    my( $condition, $offset ) = @_;

    if( not defined $offset ) { $offset = $condition; $condition = ''; }

    if( !$condition ) {
	print "      z80.jr()\n";
    } else {
	my $condition_string;
	if( defined $not{$condition} ) {
	    $condition_string = "(z80.f & FLAG_$flag{$condition}) == 0";
	} else {
	    $condition_string = "(z80.f & FLAG_$flag{$condition}) != 0";
	}
	print << "JR";
      if( $condition_string ) {
        z80.jr()
      } else {
        z80.memory.contendRead( z80.pc, 3 );
      }
JR
    }

    print "      z80.pc++;\n";
}

sub opcode_LD (@) {

    my( $dest, $src ) = @_;

    if( length $dest == 1 or $dest =~ /^REGISTER[HL]$/ ) {

	if( length $src == 1 or $src =~ /^REGISTER[HL]$/ ) {

	    if( $dest eq 'R' and $src eq 'A' ) {
		print << "LD";
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      /* Keep the RZX instruction counter right */
      rzxInstructionsOffset += ( int(z80.r) - int(z80.a))
      z80.r, z80.r7 = uint16(z80.a), uint16(z80.a)
LD
            } elsif( $dest eq 'A' and $src eq 'R' ) {
		print << "LD";
      z80.memory.contendReadNoMreq( z80.IR(), 1 );
      z80.a = byte((z80.r & 0x7f) | (z80.r7 & 0x80))
      z80.f = ( z80.f & FLAG_C ) | z80.sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)
LD
	    } else {
		print "      z80.memory.contendReadNoMreq( z80.IR(), 1 );\n" if $src eq 'I' or $dest eq 'I';
                my($lcdest, $lcsrc);
		$lcdest = lc($dest); $lcsrc = lc($src);
		print "      z80.$lcdest = z80.$lcsrc;\n" if $dest ne $src;
		if( $dest eq 'A' and $src eq 'I' ) {
		    print "      z80.f = ( z80.f & FLAG_C ) | z80.sz53Table[z80.a] | ternOpB(z80.iff2 != 0, FLAG_V, 0)\n";
		}
	    }
	} elsif( $src eq 'nn' ) {
            $dest = lc($dest);
	    print "      z80.$dest = z80.memory.readByte(z80.pc);\n      z80.pc++\n";
	} elsif( $src =~ /^\(..\)$/ ) {
	    my $register = substr $src, 1, 2;
	    $dest = lc($dest);
	    $src = lc($src);
	    print << "LD";
      z80.$dest = z80.memory.readByte(z80.$register());
LD
        } elsif( $src eq '(nnnn)' ) {
	    print << "LD";
	var wordtemp uint16
	wordtemp = uint16(z80.memory.readByte(z80.pc))
	z80.pc++
	wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
	z80.pc++
	z80.a = z80.memory.readByte(wordtemp)
LD
        } elsif( $src eq '(REGISTER+dd)' ) {
            $dest = lc($dest);
	    print << "LD";
	var offset byte;
	offset = z80.memory.readByte( z80.pc );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++;
	z80.$dest = z80.memory.readByte(uint16(int(z80.REGISTER()) + int(signExtend(offset))))
LD
        }

    } elsif( length $dest == 2 or $dest eq 'REGISTER' ) {

	my( $high, $low );

	if( $dest eq 'SP' or $dest eq 'REGISTER' ) {
	    ( $high, $low ) = ( "${dest}H", "${dest}L" );
	} else {
	    ( $high, $low ) = ( $dest =~ /^(.)(.)$/ );
	    # $low = "z80.".lc($low);
	    # $high = "z80.".lc($high);
	}

	if( $src eq 'nnnn' ) {

	    print << "LD";
      b1 := z80.memory.readByte(z80.pc)
      z80.pc++
      b2 := z80.memory.readByte(z80.pc)
      z80.pc++
      z80.set$high$low(joinBytes(b2, b1))
LD
        } elsif( $src eq 'HL' or $src eq 'REGISTER' ) {
	    print << "LD";
      z80.memory.contendReadNoMreq( z80.IR(), 1 )
      z80.memory.contendReadNoMreq( z80.IR(), 1 )
      z80.sp = z80.$src()
LD
        } elsif( $src eq '(nnnn)' ) {
	    $low = lc($low); $high = lc($high);
            if($low eq 'spl') {
		print "      sph, spl := splitWord(z80.sp)\nz80.ld16rrnn(&spl, &sph)\nz80.sp = joinBytes(sph, spl)\nbreak\n";
	    } else {
		print "      z80.ld16rrnn(&z80.$low, &z80.$high)\nbreak\n";
	    }
	}

    } elsif( $dest =~ /^\(..\)$/ ) {

	my $register = substr $dest, 1, 2;

	if( length $src == 1 ) {
	    $src = lc($src);
	    print << "LD";
      z80.memory.writeByte(z80.$register(),z80.$src);
LD
	} elsif( $src eq 'nn' ) {
	    print << "LD";
      z80.memory.writeByte(z80.$register(),z80.memory.readByte(z80.pc))
      z80.pc++
LD
        }

    } elsif( $dest eq '(nnnn)' ) {

	if( $src eq 'A' ) {
	    print << "LD";
	var wordtemp uint16 = uint16(z80.memory.readByte(z80.pc))
        z80.pc++
	wordtemp |= uint16(z80.memory.readByte(z80.pc)) << 8
        z80.pc++
	z80.memory.writeByte(wordtemp, z80.a)
LD
        } elsif( $src =~ /^(.)(.)$/ or $src eq 'REGISTER' ) {

	    my( $high, $low );

	    if( $src eq 'SP' or $src eq 'REGISTER' ) {
		( $high, $low ) = ( "${src}H", "${src}L" );
	    } else {
		( $high, $low ) = ( $1, $2 );
	    }
	    $low = lc($low); $high = lc($high);
            if($low eq 'spl') {
		print "      sph, spl := splitWord(z80.sp)\nz80.ld16nnrr(spl, sph)\nbreak\n";
	    } else {	       
		print "      z80.ld16nnrr(z80.$low, z80.$high)\nbreak\n";
	    }
	}
    } elsif( $dest eq '(REGISTER+dd)' ) {
	$src = lc($src);
	if( length $src == 1 ) {
	print << "LD";
	offset := z80.memory.readByte( z80.pc )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++;
	z80.memory.writeByte(uint16(int(z80.REGISTER()) + int(signExtend(offset))), z80.$src );
LD
        } elsif( $src eq 'nn' ) {
	    print << "LD";
	offset := z80.memory.readByte( z80.pc )
        z80.pc++
	value := z80.memory.readByte( z80.pc );
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++;
	z80.memory.writeByte(uint16(int(z80.REGISTER()) + int(signExtend(offset))), value );
LD
        }
    }

}

sub opcode_LDD (@) { ldi_ldd( 'LDD' ); }

sub opcode_LDDR (@) { ldir_lddr( 'LDDR' ); }

sub opcode_LDI (@) { ldi_ldd( 'LDI' ); }

sub opcode_LDIR (@) { ldir_lddr( 'LDIR' ); }

sub opcode_NEG (@) {
    print << "NEG";
	bytetemp := z80.a
	z80.a = 0
	z80.sub(bytetemp)
NEG
}

sub opcode_NOP (@) { }

sub opcode_OR (@) { arithmetic_logical( 'OR', $_[0], $_[1] ); }

sub opcode_OTDR (@) { otir_otdr( 'OTDR' ); }

sub opcode_OTIR (@) { otir_otdr( 'OTIR' ); }

sub opcode_OUT (@) {

    my( $port, $register ) = @_;

    if( $port eq '(nn)' and $register eq 'A' ) {
	print << "OUT";
	var outtemp uint16
	outtemp = uint16(z80.memory.readByte(z80.pc)) + (uint16(z80.a) << 8)
        z80.pc++
	z80.writePort(outtemp, z80.a)
OUT
    } elsif( $port eq '(C)' and length $register == 1 ) {
	if($register eq 0) { 
	    print "      z80.writePort(z80.BC(), $register );\n";
	} else {
	    my $lcregister = lc($register);
	    print "      z80.writePort(z80.BC(), z80.$lcregister );\n";
	}
    }
}


sub opcode_OUTD (@) { outi_outd( 'OUTD' ); }

sub opcode_OUTI (@) { outi_outd( 'OUTI' ); }

sub opcode_POP (@) { push_pop( 'POP', $_[0] ); }

sub opcode_PUSH (@) {

    my( $regpair ) = @_;

    print "      z80.memory.contendReadNoMreq( z80.IR(), 1 );\n";
    push_pop( 'PUSH', $regpair );
}

sub opcode_RES (@) { res_set( 'RES', $_[0], $_[1] ); }

sub opcode_RET (@) {

    my( $condition ) = @_;

    if( not defined $condition ) {
	print "      z80.ret();\n";
    } else {
	print "      z80.memory.contendReadNoMreq( z80.IR(), 1 );\n";
	
	if( $condition eq 'NZ' ) {
	    print << "RET";
      if( z80.pc==0x056c || z80.pc == 0x0112 ) {
	  if(z80.tapeLoadTrap() == 0) { break }
      }
RET
        }

	if( defined $not{$condition} ) {
	    print "      if(!((z80.f & FLAG_$flag{$condition}) != 0)) { z80.ret() }\n";
	} else {
	    print "      if((z80.f & FLAG_$flag{$condition}) != 0) { z80.ret() }\n";
	}
    }
}

sub opcode_RETN (@) { 

    print << "RETN";
      z80.iff1 = z80.iff2
      z80.ret()
RETN
}

sub opcode_RL (@) { rotate_shift( 'RL', $_[0] ); }

sub opcode_RLC (@) { rotate_shift( 'RLC', $_[0] ); }

sub opcode_RLCA (@) {
    print << "RLCA";
      z80.a = ( z80.a << 1 ) | ( z80.a >> 7 );
      z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |
	( z80.a & ( FLAG_C | FLAG_3 | FLAG_5 ) );
RLCA
}

sub opcode_RLA (@) {
    print << "RLA";
	var bytetemp byte = z80.a;
	z80.a = ( z80.a << 1 ) | ( z80.f & FLAG_C )
	z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & ( FLAG_3 | FLAG_5 ) ) | ( bytetemp >> 7 )
RLA
}

sub opcode_RLD (@) {
    print << "RLD";
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.writeByte(z80.HL(), (bytetemp << 4 ) | ( z80.a & 0x0f ) )
	z80.a = ( z80.a & 0xf0 ) | ( bytetemp >> 4 )
	z80.f = ( z80.f & FLAG_C ) | z80.sz53pTable[z80.a]
RLD
}

sub opcode_RR (@) { rotate_shift( 'RR', $_[0] ); }

sub opcode_RRA (@) {
    print << "RRA";
	var bytetemp byte = z80.a
	z80.a = ( z80.a >> 1 ) | ( z80.f << 7 )
	z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & ( FLAG_3 | FLAG_5 ) ) | ( bytetemp & FLAG_C ) ;
RRA
}

sub opcode_RRC (@) { rotate_shift( 'RRC', $_[0] ); }

sub opcode_RRCA (@) {
    print << "RRCA";
      z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) | ( z80.a & FLAG_C );
      z80.a = ( z80.a >> 1) | ( z80.a << 7 );
      z80.f |= ( z80.a & ( FLAG_3 | FLAG_5 ) );
RRCA
}

sub opcode_RRD (@) {
    print << "RRD";
	var bytetemp byte = z80.memory.readByte( z80.HL() )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.contendReadNoMreq( z80.HL(), 1 ); z80.memory.contendReadNoMreq( z80.HL(), 1 )
	z80.memory.writeByte(z80.HL(),  ( z80.a << 4 ) | ( bytetemp >> 4 ) )
	z80.a = ( z80.a & 0xf0 ) | ( bytetemp & 0x0f )
	z80.f = ( z80.f & FLAG_C ) | z80.sz53pTable[z80.a]
RRD
}

sub opcode_RST (@) {

    my( $value ) = @_;

    printf "      z80.memory.contendReadNoMreq( z80.IR(), 1 );\n      z80.rst(0x%02x);\n", hex $value;
}

sub opcode_SBC (@) { arithmetic_logical( 'SBC', $_[0], $_[1] ); }

sub opcode_SCF (@) {
    print << "SCF";
      z80.f = ( z80.f & ( FLAG_P | FLAG_Z | FLAG_S ) ) |
	  ( z80.a & ( FLAG_3 | FLAG_5          ) ) |
	  FLAG_C;
SCF
}

sub opcode_SET (@) { res_set( 'SET', $_[0], $_[1] ); }

sub opcode_SLA (@) { rotate_shift( 'SLA', $_[0] ); }

sub opcode_SLL (@) { rotate_shift( 'SLL', $_[0] ); }

sub opcode_SRA (@) { rotate_shift( 'SRA', $_[0] ); }

sub opcode_SRL (@) { rotate_shift( 'SRL', $_[0] ); }

sub opcode_SUB (@) { arithmetic_logical( 'SUB', $_[0], $_[1] ); }

sub opcode_XOR (@) { arithmetic_logical( 'XOR', $_[0], $_[1] ); }

sub opcode_slttrap ($) {
    print "      z80.sltTrap(int16(z80.HL()), z80.a)\n";
}

sub opcode_shift (@) {

    my( $opcode ) = @_;

    my $lc_opcode = lc $opcode;

    if( $opcode eq 'DDFDCB' ) {

	print << "shift";
      
	var tempaddr uint16
	var opcode3 byte
	z80.memory.contendRead( z80.pc, 3 )
	tempaddr = uint16(int(z80.REGISTER()) + int(signExtend(z80.memory.readByteInternal( z80.pc ))))
	z80.pc++; z80.memory.contendRead( z80.pc, 3 )
	opcode3 = z80.memory.readByteInternal( z80.pc )
	z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.memory.contendReadNoMreq( z80.pc, 1 ); z80.pc++

	switch(opcode3) {
        <%= opcodes_ddfdcb %>
	}
      
shift
    } else {
	print << "shift";
      {
	var opcode2 byte
	z80.memory.contendRead( z80.pc, 4 )
	opcode2 = z80.memory.readByteInternal( z80.pc ); z80.pc++
	z80.r++

	switch(opcode2) {
shift

    if( $opcode eq 'DD' or $opcode eq 'FD' ) {
	my $register = ( $opcode eq 'DD' ? 'IX' : 'IY' );
	my $lcregister = lc($register);
	print << "shift";
        <% register = "$lcregister" %>
        <%= opcodes_ddfd.\
            gsub(/REGISTER/i, register).\
            gsub("ix()", "IX()").\
            gsub("setixHixL", "setIX").\
            gsub("incixH", "incIXH").\
            gsub("decixH", "decIXH").
            gsub("incixL", "incIXL").\
            gsub("decixL", "decIXL").\
            gsub("z80.ix()", "z80.IX()").\
            gsub("ixH", "z80.IXH()").\
            gsub("ixL", "z80.IXL()").\
            gsub("iy()", "IY()").\
            gsub("setiyHiyL", "setIY").\
            gsub("inciyH", "incIYH").\
            gsub("deciyH", "decIYH").
            gsub("inciyL", "incIYL").\
            gsub("deciyL", "decIYL").\
            gsub("z80.iy()", "z80.IY()").\
            gsub("iyH", "z80.IYL()").\
            gsub("iyL", "z80.IYH()")
        %>
shift
        } elsif( $opcode eq 'CB' or $opcode eq 'ED' ) {
	    print "<%= opcodes_$lc_opcode %>\n";
        }

        print << "shift"
	}
      }
shift
    }
}

# Description of each file

my %description = (

    'opcodes_cb.dat'     => 'z80_cb.c: Z80 CBxx opcodes',
    'opcodes_ddfd.dat'   => 'z80_ddfd.c Z80 {DD,FD}xx opcodes',
    'opcodes_ddfdcb.dat' => 'z80_ddfdcb.c Z80 {DD,FD}CBxx opcodes',
    'opcodes_ed.dat'     => 'z80_ed.c: Z80 CBxx opcodes',
    'opcodes_base.dat'   => 'opcodes_base.c: unshifted Z80 opcodes',

);

# Main program

( my $data_file = $ARGV[0] ) =~ s!.*/!!;

print Fuse::GPL( $description{ $data_file }, '1999-2003 Philip Kendall' );

print << "COMMENT";

/* NB: this file is autogenerated by '$0' from '$data_file',
   and included in 'z80_ops.c' */

COMMENT

while(<>) {

    # Remove comments
    s/#.*//;

    # Skip (now) blank lines
    next if /^\s*$/;

    chomp;

    my( $number, $opcode, $arguments, $extra ) = split;

    if( not defined $opcode ) {
	print "    case $number: fallthrough \n";
	next;
    }

    $arguments = '' if not defined $arguments;
    my @arguments = split ',', $arguments;

    print "    case $number:\t\t/* $opcode";

    print ' ', join ',', @arguments if @arguments;
    print " $extra" if defined $extra;

    print " */\n";

    # Handle the undocumented rotate-shift-or-bit and store-in-register
    # opcodes specially

    if( defined $extra ) {

	my( $register, $opcode ) = @arguments;
	my $lcregister = lc($register);
	my $lcopcode = lc($opcode);

	if( $opcode eq 'RES' or $opcode eq 'SET' ) {

	    my( $bit ) = split ',', $extra;

	    my $operator = ( $opcode eq 'RES' ? '&' : '|' );
	    my $hexmask = res_set_hexmask( $opcode, $bit );

	    print << "CODE";
      z80.$lcregister = z80.memory.readByte(tempaddr) $operator $hexmask
      z80.memory.contendReadNoMreq(tempaddr, 1 )
      z80.memory.writeByte(tempaddr, z80.$lcregister)
      break
CODE
	} else {

	    print << "CODE";
      z80.$lcregister = z80.memory.readByte(tempaddr)
      z80.memory.contendReadNoMreq( tempaddr, 1 )
      z80.$lcopcode(&z80.$lcregister)
      z80.memory.writeByte(tempaddr, z80.$lcregister)
      break
CODE
	}
	next;
    }

    {
	no strict qw( refs );

	if( defined &{ "opcode_$opcode" } ) {
	    "opcode_$opcode"->( @arguments );
	}
    }

    print "      break\n";
}

if( $data_file eq 'opcodes_ddfd.dat' ) {

    print << "CODE";
    default:		/* Instruction did not involve H or L, so backtrack
			   one instruction and parse again */
      z80.pc--
      z80.r--
      opcode = opcode2;

      goto EndOpcode
CODE

} elsif( $data_file eq 'opcodes_ed.dat' ) {
    print << "NOPD";
    default:		/* All other opcodes are NOPD */
      break;
NOPD
}
