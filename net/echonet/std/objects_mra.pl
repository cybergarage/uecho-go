#!/usr/bin/perl
# Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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
// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

func newStandardObject(clsName string, grpCode byte, clsCode byte) Object {
	obj := NewObject()
	obj.SetClassName(clsName)
	obj.SetClassGroupCode(grpCode)
	obj.SetClassCode(clsCode)
	return obj
}

func newStandardProperty(code PropertyCode, name string, dataType string, dataSize int, getRule string, setRule string, annoRule string) Property {
	strAttrToPropertyAttr := func(strAttr string) PropertyAttribute {
		switch strAttr {
		case "required":
			return Required
		case "optional":
			return Optional
		}
		return Prohibited
	}
	prop := NewProperty(
		WithPropertyCode(code),
		WithPropertyName(name),
		WithPropertyReadAttribute(strAttrToPropertyAttr(getRule)),
		WithPropertyWriteAttribute(strAttrToPropertyAttr(setRule)),
		WithPropertyAnnoAttribute(strAttrToPropertyAttr(annoRule)),
	)
	return prop
}

// nolint:misspell, whitespace, maintidx
func (db *StandardDatabase) initObjects() {
  var obj maObject

HEADER

my $mra_definitions_file = $mra_root_dir . "/mraData/definitions/definitions.json";
open(DEF_JSON_FILE, $mra_definitions_file) or die "Failed to open $mra_definitions_file: $!";
my $def_json_data = join('',<DEF_JSON_FILE>);
close(DEF_JSON_FILE);
my $def_json = decode_json($def_json_data);
my $def_json_root = %{$def_json}{'definitions'};

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
    my $data = %{$prop}{'data'};
    my $data_type = %{$data}{'type'};
    my $data_size = %{$data}{'size'};
    my $data_ref = %{$data}{'$ref'};
    if (0< length($data_ref)) {
      my @data_refs = split(/\//, $data_ref);
      my $data_ref_len = @data_refs;
      my $data_ref_id = $data_refs[$data_ref_len -1];
      my $prop_def = %{$def_json_root}{$data_ref_id};
      $data_type = %{$prop_def}{'type'};  
      $data_size = %{$prop_def}{'size'};
      my $enums = %{$prop_def}{'enum'};
      if (0 < @data_refs) {
        foreach $enum(@{$enums}) {
          my $edt = %{$enum}{'edt'};
          my $name = %{$enum}{'name'};
          my $descs = %{$enum}{'descriptions'};
          my $desc = %{$descs}{'en'};
          if (0< length($edt)) {
          }
        }
      }
    }
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