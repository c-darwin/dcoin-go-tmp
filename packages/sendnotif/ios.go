// +build darwin
// +build arm arm64

package sendnotif

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework GLKit -framework UIKit
#import <UIKit/UIKit.h>
#import <Foundation/Foundation.h>
#import <GLKit/GLKit.h>

void ShowMess(char* title, text) {
	UIAlertView* alert = [[UIAlertView alloc] initWithTitle:@title message:@text delegate:nil cancelButtonTitle:@"OK" otherButtonTitles: nil];
	[alert show];
	[alert release];
}

*/
import "C"

func SendMobileNotification(title, text string) {
	C.ShowMess();
}