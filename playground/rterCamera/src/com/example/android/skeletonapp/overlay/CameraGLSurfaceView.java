package com.example.android.skeletonapp.overlay;

import android.annotation.TargetApi;
import android.content.Context;
import android.graphics.PixelFormat;
import android.opengl.GLSurfaceView;
import android.os.Build;

@TargetApi(Build.VERSION_CODES.GINGERBREAD)
public class CameraGLSurfaceView extends GLSurfaceView {

	public CameraGLSurfaceView(Context context) {
		super(context);
		
		//needed to overlay gl view over camera preview
		this.setZOrderMediaOverlay(true);
		
        // Create an OpenGL ES 1.0 context
        this.setEGLContextClientVersion(1);
        
        this.getHolder().setFormat(PixelFormat.TRANSLUCENT);
        this.setEGLConfigChooser(8, 8, 8, 8, 16, 0);
        
        // Set the Renderer for drawing on the GLSurfaceView
        this.setRenderer(new CameraGLRenderer(context));
           
        // Render the view only when there is a change in the drawing data
        //this.setRenderMode(GLSurfaceView.RENDERMODE_WHEN_DIRTY);
	}

}
