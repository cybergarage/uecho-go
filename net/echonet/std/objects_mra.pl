#!/usr/bin/perl
# Copyright (C) 2018 Satoshi Konno. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


use utf8;
use JSON;
use File::Find;

if (@ARGV < 1){
  exit 1;
}
my $mra_root_dir = $ARGV[0];

print<<HEADER;
// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

func newStandardObject(clsName string, grpCode byte, clsCode byte) *Object {
	obj := NewObject()
	obj.SetClassName(clsName)
	obj.SetClassGroupCode(grpCode)
	obj.SetClassCode(clsCode)
	return obj
}

func newStandardProperty(code PropertyCode, name string, dataType string, dataSize int, getRule string, setRule string, annoRule string) *Property {
	strAttrToPropertyAttr := func(strAttr string) PropertyAttr {
		switch strAttr {
		case "required":
			return Required
		case "optional":
			return Optional
		}
		return Prohibited
	}
	prop := NewProperty()
	prop.SetCode(code)
	prop.SetName(name)
	prop.setReadAttribute(strAttrToPropertyAttr(getRule))
	prop.setWriteAttribute(strAttrToPropertyAttr(setRule))
	prop.setAnnoAttribute(strAttrToPropertyAttr(annoRule))
	return prop
}

// nolint:misspell, whitespace, maintidx
func (db *StandardDatabase) initObjects() {
  var obj *Object

HEADER

my @mra_sub_dirs = (
  "/mraData/superClass/",
  "/mraData/nodeProfile/",
  "/mraData/devices/"
);

my @device_json_files;
foreach my $mra_sub_dir(@mra_sub_dirs){
  my $mra_root_dir = $mra_root_dir . $mra_sub_dir;
  find sub {
      my $file = $_;
      my $path = $File::Find::name;
      if(-f $file){
        push(@device_json_files, $path);
      }
  }, $mra_root_dir;
}

foreach my $device_json_file(@device_json_files){
  open(DEV_JSON_FILE, $device_json_file) or die "$!";
  my $device_json_data = join('',<DEV_JSON_FILE>);
  close(DEV_JSON_FILE);
  my $device_json = decode_json($device_json_data);

  my $cls_names = %{$device_json}{'className'};
  my $cls_name = %{$cls_names}{'en'};
  my $grp_cls_code = %{$device_json}{'eoj'};
  my $grp_code = substr($grp_cls_code, 2, 2);
  my $cls_code = substr($grp_cls_code, 4);
  printf("// %s (0x%s%s)\n", $cls_name, $grp_code, $cls_code);
  printf("obj = newStandardObject(\"%s\", 0x%s, 0x%s)\n", $cls_name, $grp_code, $cls_code);

  my $props = %{$device_json}{'elProperties'};
  foreach $prop(@{$props}) {
    my $epc = %{$prop}{'epc'};
    my $names = %{$prop}{'propertyName'};
    my $name = %{$names}{'en'};
    my $rules = %{$prop}{'accessRule'};
    my $get_rule = %{$rules}{'get'};
    my $set_rule = %{$rules}{'set'};
    my $anno_rule = %{$rules}{'inf'};
    my $data_type = "";
    printf("obj.AddProperty(newStandardProperty(%s, \"%s\", \"%s\", %d, \"%s\", \"%s\", \"%s\"))\n",
      $epc,
      $name,
      $data_type,
      $data_size,
      $get_rule,
      $set_rule,
      $anno_rule,
      );
   }
  printf("db.addObject(obj)\n\n", $grp_code, $cls_code);
}
print<<FOTTER;
}
FOTTER