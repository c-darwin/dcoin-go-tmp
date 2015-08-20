package org.golang.app;

import android.app.Activity;
import android.app.NativeActivity;
import android.content.Context;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.util.Log;
import android.content.Intent;
import android.content.ComponentName;
import android.net.Uri;

public class GoNativeActivity extends NativeActivity {

    private static GoNativeActivity goNativeActivity;

    public GoNativeActivity() {
        super();
        goNativeActivity = this;
    }

    String getTmpdir() {
        return getCacheDir().getAbsolutePath();
    }

    String getFilesdir() {
        return getExternalFilesDir(null).getAbsolutePath();
    }

    private void load() {

        // Interestingly, NativeActivity uses a different method
        // to find native code to execute, avoiding
        // System.loadLibrary. The result is Java methods
        // implemented in C with JNIEXPORT (and JNI_OnLoad) are not
        // available unless an explicit call to System.loadLibrary
        // is done. So we do it here, borrowing the name of the
        // library from the same AndroidManifest.xml metadata used
        // by NativeActivity.
        try {
            ActivityInfo ai = getPackageManager().getActivityInfo(
                    getIntent().getComponent(), PackageManager.GET_META_DATA);
            if (ai.metaData == null) {
                Log.e("Go", "loadLibrary: no manifest metadata found");
                return;
            }
            String libName = ai.metaData.getString("android.app.lib_name");
            System.loadLibrary(libName);
        } catch (Exception e) {
            Log.e("Go", "loadLibrary failed", e);
        }

        try {
            Intent intent = new Intent(Intent.ACTION_VIEW);
            Uri data = Uri.parse("http://localhost:8089");
            intent.setData(data);
            startActivity(intent);
        } catch (Exception e) {
            Log.e("Go", "http://localhost:8089 failed", e);
        }
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        load();
        super.onCreate(savedInstanceState);
    }
}
