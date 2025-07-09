import React, { useRef, useEffect } from 'react';
import { Box, Typography } from '@mui/material';
import Hls from 'hls.js';

interface VideoPlayerProps {
  url: string;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ url }) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [error, setError] = React.useState<string | null>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    let hls: Hls | null = null;

    const setupHls = () => {
      if (Hls.isSupported()) {
        hls = new Hls({
          enableWorker: true,
          lowLatencyMode: true,
          backBufferLength: 90
        });
        
        hls.loadSource(url);
        hls.attachMedia(video);
        
        hls.on(Hls.Events.MANIFEST_PARSED, () => {
          video.play().catch(err => {
            console.error('Error playing video:', err);
            setError('Failed to autoplay video. Please click to play.');
          });
        });
        
        hls.on(Hls.Events.ERROR, (_event, data) => {
          if (data.fatal) {
            switch (data.type) {
              case Hls.ErrorTypes.NETWORK_ERROR:
                console.error('Network error:', data);
                hls?.startLoad();
                break;
              case Hls.ErrorTypes.MEDIA_ERROR:
                console.error('Media error:', data);
                hls?.recoverMediaError();
                break;
              default:
                console.error('Unrecoverable error:', data);
                hls?.destroy();
                setError('Failed to load video stream. Please try again later.');
                break;
            }
          }
        });
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        // For Safari
        video.src = url;
        video.addEventListener('loadedmetadata', () => {
          video.play().catch(err => {
            console.error('Error playing video:', err);
            setError('Failed to autoplay video. Please click to play.');
          });
        });
      } else {
        setError('Your browser does not support HLS video streaming.');
      }
    };

    setupHls();

    return () => {
      if (hls) {
        hls.destroy();
      }
      
      if (video) {
        video.pause();
        video.src = '';
        video.load();
      }
    };
  }, [url]);

  const handleClick = () => {
    if (videoRef.current && videoRef.current.paused) {
      videoRef.current.play().catch(err => {
        console.error('Error playing video on click:', err);
      });
    }
  };

  return (
    <Box 
      sx={{ 
        position: 'absolute', 
        top: 0, 
        left: 0, 
        width: '100%', 
        height: '100%',
        bgcolor: 'black',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center'
      }}
    >
      {error ? (
        <Typography color="error" variant="body1">
          {error}
        </Typography>
      ) : (
        <video
          ref={videoRef}
          style={{ 
            width: '100%', 
            height: '100%', 
            objectFit: 'contain' 
          }}
          playsInline
          onClick={handleClick}
          controls
        />
      )}
    </Box>
  );
};

export default VideoPlayer;
