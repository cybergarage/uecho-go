// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechosearch is a search utility for Echonet Lite.

	NAME
	uechosearch

	SYNOPSIS
	uechosearch [OPTIONS]

	DESCRIPTION
	uechosearch is a search utility for Echonet Lite.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet"
)

void uecho_search_print_messages(uEchoController *ctrl, uEchoMessage *msg)
{
  uEchoProperty *prop;
  size_t opc, n;
  
  opc = uecho_message_getopc(msg);
  printf("%s %1X %1X %02X %03X %03X %02X %ld ",
         uecho_message_getsourceaddress(msg),
         uecho_message_getehd1(msg),
         uecho_message_getehd2(msg),
         uecho_message_gettid(msg),
         uecho_message_getsourceobjectcode(msg),
         uecho_message_getdestinationobjectcode(msg),
         uecho_message_getesv(msg),
         opc);
  
  for (n=0; n<opc; n++) {
    prop = uecho_message_getproperty(msg, n);
    printf("%02X", uecho_property_getcode(prop));
  }
  
  printf("\n");
}

func main() {
	flag.Bool("v", false, "verbose")
	flag.Parse()

	ctrl := echonet.NewController()

	err := ctrl.Start()
	if err != nil {
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	for _, node := range ctrl.GetNodes() {
		objs := node.GetObjects()
		if len(objs) <= 0 {
			fmt.Printf("%s\n", node.GetAddress())
			continue
		}
		for _, obj := range objs {
			fmt.Printf("%s %06X\n", node.GetAddress(), obj.GetCode())
		}
	}

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
