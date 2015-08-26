package org.golang.app;

import android.app.Activity;
import android.os.Bundle;
import android.content.Intent;
import android.util.Log;
import android.net.Uri;
import android.os.Handler;
import android.os.SystemClock;

public class MainActivity extends Activity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
	  
	/*Intent dialogIntent = new Intent("org.golang.app.GoNativeActivity");
	dialogIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
	this.startActivity(dialogIntent);*/
	
	
         super.onCreate(savedInstanceState);  
	  Intent intent=new Intent("org.golang.app.MyService");  
	  this.startService(intent);
	  
	  
	    SystemClock.sleep(3000);
	    
	   /* Handler handler = new Handler(); 
	    handler.postDelayed(new Runnable() { 
		public void run() {*/
		  try {
			  Intent intent1 = new Intent(Intent.ACTION_VIEW);
			  Uri data = Uri.parse("http://localhost:8089");
			  intent1.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
			  intent1.setData(data);
			  startActivity(intent1);
		  } catch (Exception e) {
			  Log.e("Go", "http://localhost:8089 failed", e);
		  }
		//} 
	    /*}, 3000); */
    

		
    }
    
    
    protected void onStart(Bundle savedInstanceState) {

		  try {
			  Intent intent1 = new Intent(Intent.ACTION_VIEW);
			  Uri data = Uri.parse("http://localhost:8089");
			  intent1.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
			  intent1.setData(data);
			  startActivity(intent1);
		  } catch (Exception e) {
			  Log.e("Go", "http://localhost:8089 failed", e);
		  }  
	  
    }

}
