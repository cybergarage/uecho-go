// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechodump is a dump utility for Echonet Lite.

	NAME
	uechopost

	SYNOPSIS
	uechopost [options] <address> <obj> <esv> <property (epc, pdc, edt) ...

	DESCRIPTION
	uechopost is a controller utility to send any messages to Echonet Lite nodes.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"fmt"
	"os"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAIL    = 1
)

func outputUsage() {
	fmt.Printf("Usage : uechopost [options] <address> <obj> <esv> <property (epc, pdc, edt) ...>\n")
}

func outputError(err error) {
	fmt.Printf("ERROR : %s\n", err.Error())
}

func main() {
	ctrl := NewPostController()

	argc := len(os.Args)
	if argc < 4 {
		outputUsage()
		os.Exit(EXIT_FAIL)
		return
	}

	err := ctrl.Start()
	if err != nil {
		outputError(err)
		os.Exit(EXIT_FAIL)
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		outputError(err)
		return
	}

	/*
	     // Find destination node

	     dstNodeAddr = argv[0];

	     dstNode = NULL;
	     for (n=0; n<UECHOPOST_RESPONSE_RETRY_COUNT; n++) {
	       uecho_sleep(UECHOPOST_MAX_RESPONSE_MTIME / UECHOPOST_RESPONSE_RETRY_COUNT);
	       dstNode = uecho_controller_getnodebyaddress(ctrl, dstNodeAddr);
	       if (dstNode)
	         break;
	     }

	     if (!dstNode) {
	       printf("Node (%s) is not found\n", dstNodeAddr);
	       uecho_controller_delete(ctrl);
	       return EXIT_FAILURE;
	     }

	     // Find destination object

	     sscanf(argv[1], "%x", &dstObjCode);

	     dstObj = uecho_node_getobjectbycode(dstNode, dstObjCode);

	     if (!dstNode) {
	       printf("Node (%s) doesn't has the specified object (%06X)\n", dstNodeAddr, dstObjCode);
	       uecho_controller_delete(ctrl);
	       return EXIT_FAILURE;
	     }

	     // Create Message

	     msg = uecho_message_new();
	     sscanf(argv[2], "%x", &esv);
	     uecho_message_setesv(msg, esv);

	   #if defined(DEBUG)
	     printf("%s %06X %01X\n", dstNodeAddr, dstObjCode, esv);
	   #endif

	     edata = edt = argv[3];
	     edtSize = strlen(argv[3]);
	     while ((edt - edata + (2 + 2)) <= edtSize) {
	       sscanf(edt, "%02x%02x", &epc, &pdc);
	       edt += (2 + 2);

	   #if defined(DEBUG)
	       printf("[%02X] = %02X ", epc, pdc);
	   #endif

	       if (pdc == 0) {
	         uecho_message_setproperty(msg, epc, 0, NULL);
	         continue;
	       }

	       if (edtSize < (edt - edata + (pdc * 2)))
	         break;

	       propData = (byte *)malloc(pdc);
	       for (n=0; n<pdc; n++) {
	         sscanf(edt, "%02x", &edtByte);
	   #if defined(DEBUG)
	         printf("%02X", edtByte);
	   #endif
	         propData[n] = edtByte & 0xFF;
	         edt += 2;
	       }
	       uecho_message_setproperty(msg, epc, pdc, propData);
	       free(propData);
	     }
	   #if defined(DEBUG)
	     printf("\n");
	   #endif

	     // Send message

	     isResponseRequired = uecho_message_isresponserequired(msg);
	     if (isResponseRequired) {
	       resMsg = uecho_message_new();
	       if (uecho_controller_postmessage(ctrl, dstObj, msg, resMsg)) {
	         uechopost_print_objectresponse(ctrl, resMsg);
	       }
	       uecho_message_delete(resMsg);
	     }
	     else {
	       uecho_controller_sendmessage(ctrl, dstObj, msg);
	     }

	     uecho_message_delete(msg);
	*/

	err = ctrl.Stop()
	if err != nil {
		return
	}

	os.Exit(EXIT_SUCCESS)
}
