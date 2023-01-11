#!/usr/bin/perl
# Copyright (C) 2018 Satoshi Konno. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if (@ARGV < 1){
  exit 1;
}
my $manlist_filename = $ARGV[0];

print<<HEADER;
// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// nolint: maintidx
func (db *StandardDatabase) initManufactures() {
HEADER

open(MANLIST, $manlist_filename) or die "$!";
while(<MANLIST>){
  chomp($_);
  $_ =~ s/(['"].*?['"])/(my $s = $1) =~ tr|,|=|; $s/eg;
  my @row = split(/(?!"),/, $_, -1);;
  my $code = $row[0];
  if (length($code ) != 6) {
    next;
  }
  my $name = $row[1];
  $name =~ s/=/,/g;
  $name =~ s/　/ /g; # converts zenkaku spaces to spaces
  if ($name !~ /^\"/) {
    $name = "\""  . $name
  }
  if ($name !~ /\"$/) {
    $name = $name . "\"" 
  }
  printf("db.addManufacture(NewManufacture(0x%s, %s))\n", $code, $name);
}
printf("db.addManufacture(NewManufacture(0xFFFFFF, \"Experimental\"))\n");
printf("db.addManufacture(NewManufacture(0xFFFFFE, \"Undefined\"))\n");
close(MANLIST);
print<<FOTTER;
}
FOTTER