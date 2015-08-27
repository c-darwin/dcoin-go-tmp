package notif

/*
#cgo LDFLAGS: -llog -landroid
#include <android/log.h>
#include <jni.h>
#include <stdlib.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Go/notif", __VA_ARGS__)
#define LOG_INFO(...) __android_log_print(ANDROID_LOG_INFO, "Go/notif", __VA_ARGS__)

void notif_manager_init(void* java_vm, void* ctx, char* title, char* text) {
	JavaVM* vm = (JavaVM*)(java_vm);
	JNIEnv* env;
	int err;
	int attached = 0;

	err = (*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6);
	if (err != JNI_OK) {
		if (err == JNI_EDETACHED) {
			if ((*vm)->AttachCurrentThread(vm, &env, 0) != 0) {
				LOG_FATAL("cannot attach JVM");
			}
			attached = 1;
		} else {
			LOG_FATAL("GetEnv unexpected error: %d", err);
		}
	}
	
	
	//char* title = "tttt";
	//char* text = "text";
	
	jstring javaTitle = (jstring)(*env)->NewStringUTF(env, (const char *)title);
	jstring javaText = (jstring)(*env)->NewStringUTF(env, (const char *)text);
		
	jclass cls = (*env)->GetObjectClass(env, ctx);
	LOG_INFO("jclasscls %x", cls);		
	jmethodID notif1 = (*env)->GetMethodID(env, cls, "notif", "(Ljava/lang/String;Ljava/lang/String;)V");
	LOG_INFO("notif1 %x", notif1);
	(jstring)(*env)->CallObjectMethod(env, ctx, notif1, javaTitle, javaText);
		

	if (attached) {
		(*vm)->DetachCurrentThread(vm);
	}
}
*/
import "C"
import (
	"golang.org/x/mobile/internal/mobileinit"
)


func NotifInit(title string, text string) {
	ctx := mobileinit.Context{}
	C.notif_manager_init(ctx.JavaVM(), ctx.AndroidContext(), C.CString(title), C.CString(text))
}


