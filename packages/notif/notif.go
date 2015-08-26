package notif 

/*
#cgo LDFLAGS: -llog -landroid
#include <stdio.h>
#include <stdlib.h>
#include <jni.h>
#include <android/log.h>
#include <android/native_activity.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Go/notif", __VA_ARGS__)
#define LOG_INFO(...) __android_log_print(ANDROID_LOG_INFO, "Go/notif", __VA_ARGS__)

static jclass find_class(JNIEnv *env, const char *class_name) {
	jclass clazz = (*env)->FindClass(env, class_name);
	if (clazz == NULL) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find %s", class_name);
		return NULL;
	}
	return clazz;
}


static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

void notif(ANativeActivity *activity) {

	//ANativeActivity *activity;

	LOG_INFO("444444444444444444444444444: %s", "546545454");
	LOG_INFO("notifactivity: %x", activity);
	LOG_INFO("notifactivity->env: %x", activity->env);
	LOG_INFO("notifactivity->clazz: %x", activity->clazz);
		
	
		JNIEnv* env = activity->env;			
		// Note that activity->clazz is mis-named.
		JavaVM* current_vm = activity->vm;
		jobject current_ctx = activity->clazz;
		
		jclass cls = (*env)->GetObjectClass(env, current_ctx);
		LOG_INFO("jclasscls %x", cls);
		jmethodID gettmpdir = find_method(env, cls, "getTmpdir", "()Ljava/lang/String;");
		LOG_INFO("gettmpdir %x", gettmpdir);
		jstring jpath = (jstring)(*env)->CallObjectMethod(env, current_ctx, gettmpdir, NULL);
		LOG_INFO("tmpdirtmpdirtmpdirtmpdir: %s", jpath);
		const char* tmpdir = (*env)->GetStringUTFChars(env, jpath, NULL);
		LOG_INFO("tmpdirtmpdirtmpdirtmpdir: %s", tmpdir);		
		
		jmethodID notif1 = find_method(env, cls, "notif", "()Ljava/lang/String;");
		LOG_INFO("notif1 %x", notif1);
		jstring jpath11 = (jstring)(*env)->CallObjectMethod(env, current_ctx, notif1, NULL);
		LOG_INFO("notif1notif1notif1: %s", jpath11);
		const char* jpath11111 = (*env)->GetStringUTFChars(env, jpath11, NULL);
		LOG_INFO("notif1notif1notif11111111: %s", jpath11111);		
		

}
*/
import "C"
import (
	"golang.org/x/mobile/app"
	//"fmt"
	"unsafe"
)

var Activity *C.ANativeActivity


func InitNotif() {
	//fmt.Println("GactivityInitNotif:", activity);
	//fmt.Println("GactivityInitNotif:", (*C.ANativeActivity)(unsafe.Pointer(activity)));
	//ctx := mobileinit.Context{}
	//fmt.Println("CTX",ctx)
	//C.notif(ctx.JavaVM(), ctx.AndroidContext(), C.CString("11111111"))
	C.notif((*C.ANativeActivity)(unsafe.Pointer(app.Gactivity)));
}
