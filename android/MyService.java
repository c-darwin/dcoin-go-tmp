package org.golang.app;

import android.app.Notification;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.os.IBinder;
import android.util.Log;
import android.os.Binder;
import java.util.concurrent.TimeUnit;
import android.widget.Toast;
import android.os.SystemClock;


public class MyService extends Service {

    
    @Override
    public IBinder onBind(Intent intent) {
    	Log.e("Go", "MyService 222200111");

        return null;
    }

    
    public int onStartCommand(Intent intent, int flags, int startId) {
    
	Log.e("Go", "MyService 00111");
         return super.onStartCommand(intent, flags, startId);

    }
    
    public void onStart() {
	Log.e("Go", "MyService 00221");
    
    }
    
    public void onDestroy() {
	super.onDestroy();
	    Log.e("Go", "MyService 00221111");
      }
    
    
    @Override
    public void onCreate() {
        super.onCreate();
        
        sendNotif();
        
	Log.e("onCreate", "MyService111111111 01");
        //GoNativeActivity.notif();
        
	Log.e("onCreate", "MyService 01");
        GoNativeActivity.load();
        //Toast.makeText(this, "************Service Started #############+", Toast.LENGTH_LONG).show();
        // do something when the service is created
        //GoNativeActivity.St();
	Log.e("startActivity", "333333333333 01");
	    
	  
	  
      //SystemClock.sleep(4000);
	    
      Intent dialogIntent = new Intent(this, GoNativeActivity.class);
      dialogIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
      startActivity(dialogIntent);


	Log.e("startActivity", "startActivity 01");
	//GoNativeActivity.St();
    }
    
    
        void sendNotif() {
    
   	  Log.e("Go", "sendNotif 1");

	  Notification notif = new Notification(R.drawable.icon, "Text in status bar",
		  System.currentTimeMillis());
		  
	  // 3-я часть
	  Intent intent = new Intent(this, MainActivity.class);
	  //intent.putExtra(GoNativeActivity.FILE_NAME, "somefile");
	  PendingIntent pIntent = PendingIntent.getActivity(this, 0, intent, 0);

	  // 2-я часть
	  notif.setLatestEventInfo(this, "Notification's title", "Notification's text", pIntent);

	  // ставим флаг, чтобы уведомление пропало после нажатия
	  //notif.flags |= Notification.FLAG_AUTO_CANCEL;

   	  Log.e("Go", "sendNotif 2");
	  // отправляем
	  startForeground(1, notif);
	  
   	  Log.e("Go", "sendNotif 3");

    }
    
}

/*
public class MyService extends Service {

  MyBinder binder = new MyBinder();
  
    NotificationManager nm;

    
    @Override
    public void onCreate() {
	Log.e("Go", "MyService 01");
        super.onCreate();
	Log.e("Go", "MyService 02");
        nm = (NotificationManager) getSystemService(NOTIFICATION_SERVICE);
	Log.e("Go", "MyService 03");
    }

    public int onStartCommand(Intent intent, int flags, int startId) {
	Log.e("Go", "MyService 001");
        sendNotif();
        
        for (int i = 1; i<=500; i++) {
	  Log.d("SERVTEST", "i = " + i);
	  try {
	    TimeUnit.SECONDS.sleep(1);
	  } catch (InterruptedException e) {
	    e.printStackTrace();
	  }
	}
    
        try {
		//System.loadLibrary("dcoin");		
		Log.e("Go", "loadLibrary ok");
			
	} catch (Exception e) {
		Log.e("Go", "loadLibrary failed");
	}
        return super.onStartCommand(intent, flags, startId);
    }


    void sendNotif() {
    
   	  Log.e("Go", "sendNotif 1");

	  Notification notif = new Notification(R.drawable.icon, "Text in status bar",
		  System.currentTimeMillis());
		  
	  // 3-я часть
	  Intent intent = new Intent(this, GoNativeActivity.class);
	  //intent.putExtra(GoNativeActivity.FILE_NAME, "somefile");
	  PendingIntent pIntent = PendingIntent.getActivity(this, 0, intent, 0);

	  // 2-я часть
	  notif.setLatestEventInfo(this, "Notification's title", "Notification's text", pIntent);

	  // ставим флаг, чтобы уведомление пропало после нажатия
	  notif.flags |= Notification.FLAG_AUTO_CANCEL;

   	  Log.e("Go", "sendNotif 2");
	  // отправляем
	  startForeground(1, notif);
	  
   	  Log.e("Go", "sendNotif 3");

    }
    
    public IBinder onBind(Intent arg0) {
      Log.d("GO", "MyService onBind");
      return binder;
    }
    
    class MyBinder extends Binder {
      MyService getService() {
	return MyService.this;
      }
    }  
}

*/

/*
public class MyService extends Service {
  
  final String LOG_TAG = "myLogs";

  public void onCreate() {
    super.onCreate();
    Log.d(LOG_TAG, "onCreate");
  }
  
  public int onStartCommand(Intent intent, int flags, int startId) {
    Log.d(LOG_TAG, "onStartCommand");
    someTask();
    return super.onStartCommand(intent, flags, startId);
  }

  public void onDestroy() {
    super.onDestroy();
    Log.d(LOG_TAG, "onDestroy");
  }

  public IBinder onBind(Intent intent) {
    Log.d(LOG_TAG, "onBind");
    return null;
  }
  
  void someTask() {
  }
}*/