package org.golang.app;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
//import android.app.AlarmManager;
//import android.app.PendingIntent;
import android.app.NativeActivity;
import android.util.Log;
import java.util.concurrent.TimeUnit;


public class BootReceiver extends BroadcastReceiver {

    public static int Started = 0;

    
/*
    @Override
    public void onReceive(Context context, Intent intent) {
         Log.d("SERVTEST0", "SERVTEST0");
	 for (int i = 1; i<=500; i++) {
	  Log.d("onReceive", "i = " + i);
	  try {
	    TimeUnit.SECONDS.sleep(1);
	  } catch (InterruptedException e) {
	    e.printStackTrace();
	  }
	}
	
        Intent service = new Intent(context, MyService.class);
        context.startService(service);
  }*/
    @Override
    public void onReceive(Context context, Intent intent) {

    	Log.e("Go", "MyService ++++++++++++++++++");
    	
    	/*
	  
      Intent dialogIntent = new Intent(context, GoNativeActivity.class);
      dialogIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
      context.startActivity(dialogIntent);*/
      
      
        if (intent.getAction().equalsIgnoreCase(Intent.ACTION_BOOT_COMPLETED)) {
            Intent serviceIntent = new Intent(context, MyService.class);
            context.startService(serviceIntent);
        }
    }
}
/*
public class BootReceiver extends BroadcastReceiver {
@Override
public void onReceive(Context context, Intent intent) {
    AlarmManager am = (AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
    PendingIntent pi = PendingIntent.getService(context, 0, new Intent(context, MyService.class), PendingIntent.FLAG_UPDATE_CURRENT);
    am.setInexactRepeating(AlarmManager.RTC_WAKEUP, System.currentTimeMillis() + interval, interval, pi);
}}*/