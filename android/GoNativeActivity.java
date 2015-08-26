package org.golang.app;

import android.app.Notification;
import android.app.NotificationManager;
import android.content.pm.ActivityInfo;
import android.app.NativeActivity;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;
import android.content.Intent;
import android.net.Uri;
import android.app.PendingIntent;
import android.content.Context;
import android.widget.Toast;
import android.app.TaskStackBuilder;
import android.support.v4.app.NotificationCompat;



public class GoNativeActivity extends NativeActivity {

    private static GoNativeActivity goNativeActivity;

    public GoNativeActivity() {
        super();
        goNativeActivity = this;
    }
    
    public static NotificationManager nm;
    
   // private static Context context;
    
    
    public String notif() {
    

   	  Log.e("Go", "notif()OK");
   	  
	  String text="65656565656";
	  
  	  Intent intent = new Intent("org.golang.app.MainActivity");	    

  	  
   	  Log.e("Go", "thisthis"+this);
   	  
	  //Toast.makeText(this, text, Toast.LENGTH_LONG).show();

	  NotificationCompat.Builder mBuilder = new NotificationCompat.Builder(this);
	  mBuilder.setSmallIcon(R.drawable.icon);
	  mBuilder.setContentTitle("!!!!!!!!!!!++++!!!!!!");
	  mBuilder.setContentText("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA");
	          
	  Intent resultIntent = new Intent(this, MainActivity.class);
	  TaskStackBuilder stackBuilder = TaskStackBuilder.create(this);
	  stackBuilder.addParentStack(MainActivity.class);

	  // Adds the Intent that starts the Activity to the top of the stack
	  stackBuilder.addNextIntent(resultIntent);
	  PendingIntent resultPendingIntent = stackBuilder.getPendingIntent(0,PendingIntent.FLAG_UPDATE_CURRENT);
	  mBuilder.setContentIntent(resultPendingIntent);

	  NotificationManager mNotificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
    
	  // notificationID allows you to update the notification later on.
	  mNotificationManager.notify(1111, mBuilder.build());

	  /*Notification notif = new Notification(R.drawable.icon, "Text in status bar",
		  System.currentTimeMillis());
		  
	  
	  Intent intent = new Intent("org.golang.app.MainActivity");	    
	  

   	  Log.d("Go", "notif()OK %v"+GoApp.getAppContext());
	  PendingIntent pIntent = PendingIntent.getActivity(GoApp.getAppContext(), 0, intent, 0);

	  /*
	  // 2-я часть
	  notif.setLatestEventInfo(GoNativeActivity.context, "title "+text, text, pIntent);

	  // ставим флаг, чтобы уведомление пропало после нажатия
	  notif.flags |= Notification.FLAG_AUTO_CANCEL;

   	  Log.e("Go", "sendNotif 2");
	  // отправляем
	  GoNativeActivity.nm.notify(1, notif);*/
   	  Log.e("Go", "sendNotif ok");
	 return "65888888888888888888888";
    }
    

    String getTmpdir() {
        return getCacheDir().getAbsolutePath();
    }

    String getFilesdir() {
        return getExternalFilesDir(null).getAbsolutePath();
    }

    public static void load() {

        // Interestingly, NativeActivity uses a different method
        // to find native code to execute, avoiding
        // System.loadLibrary. The result is Java methods
        // implemented in C with JNIEXPORT (and JNI_OnLoad) are not
        // available unless an explicit call to System.loadLibrary
        // is done. So we do it here, borrowing the name of the
        // library from the same AndroidManifest.xml metadata used
        // by NativeActivity.
		try {
			Log.d("Gomylib", "6666666666666");
			
			System.loadLibrary("dcoin");
			Log.d("Gomylib", "77777777777777777");
			
		
			
		} catch (Exception e) {
			Log.e("Go", "loadLibrary failed", e);
		}
		



    }

    public void onStart(Bundle savedInstanceState) {
    
    		  try {
			  Intent intent1 = new Intent(Intent.ACTION_VIEW);
			  Uri data = Uri.parse("http://localhost:8089");
			  intent1.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
			  intent1.setData(data);
			  startActivity(intent1);
		  } catch (Exception e) {
			  Log.e("Go", "http://localhost:8089 failed", e);
		  }
	Log.d("Go1111", "startService 00111+");
    }
    @Override
    public void onCreate(Bundle savedInstanceState) {
    
	Log.d("Go1111", "startService 001+");
	
        //load();
        
        moveTaskToBack(true);
        
        /*
          super.onCreate(savedInstanceState);  
	  Intent intent=new Intent("org.golang.app.MyService");  
	  this.startService(intent);
      
	
	try {		
	    startService(new Intent("org.golang.app.MyService"));
        } catch (Exception e) {
	    Log.e("Goerr", "errr", e);
	}
	Log.d("Go", "startService 0001+");

        */

        super.onCreate(savedInstanceState);
        
        
        //GoNativeActivity.context = getApplicationContext();
        
    }
}
