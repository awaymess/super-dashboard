'use client';

import { useEffect, useRef } from 'react';

export function LiquidBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let animationFrameId: number;
    let time = 0;

    const resize = () => {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
    };

    const drawBlob = (x: number, y: number, radius: number, color: string, offset: number) => {
      ctx.beginPath();
      for (let i = 0; i < 360; i++) {
        const angle = (i * Math.PI) / 180;
        const wobble = Math.sin(angle * 3 + time + offset) * 20;
        const r = radius + wobble;
        const px = x + Math.cos(angle) * r;
        const py = y + Math.sin(angle) * r;
        if (i === 0) {
          ctx.moveTo(px, py);
        } else {
          ctx.lineTo(px, py);
        }
      }
      ctx.closePath();
      ctx.fillStyle = color;
      ctx.fill();
    };

    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      
      // Create gradient blobs
      const gradient1 = ctx.createRadialGradient(
        canvas.width * 0.3 + Math.sin(time * 0.5) * 50,
        canvas.height * 0.4 + Math.cos(time * 0.3) * 50,
        0,
        canvas.width * 0.3,
        canvas.height * 0.4,
        300
      );
      gradient1.addColorStop(0, 'rgba(59, 130, 246, 0.3)');
      gradient1.addColorStop(1, 'rgba(59, 130, 246, 0)');

      const gradient2 = ctx.createRadialGradient(
        canvas.width * 0.7 + Math.cos(time * 0.4) * 50,
        canvas.height * 0.6 + Math.sin(time * 0.6) * 50,
        0,
        canvas.width * 0.7,
        canvas.height * 0.6,
        300
      );
      gradient2.addColorStop(0, 'rgba(139, 92, 246, 0.3)');
      gradient2.addColorStop(1, 'rgba(139, 92, 246, 0)');

      drawBlob(
        canvas.width * 0.3 + Math.sin(time * 0.5) * 50,
        canvas.height * 0.4 + Math.cos(time * 0.3) * 50,
        200,
        gradient1 as unknown as string,
        0
      );

      drawBlob(
        canvas.width * 0.7 + Math.cos(time * 0.4) * 50,
        canvas.height * 0.6 + Math.sin(time * 0.6) * 50,
        250,
        gradient2 as unknown as string,
        Math.PI
      );

      time += 0.01;
      animationFrameId = requestAnimationFrame(animate);
    };

    resize();
    animate();
    window.addEventListener('resize', resize);

    return () => {
      window.removeEventListener('resize', resize);
      cancelAnimationFrame(animationFrameId);
    };
  }, []);

  return (
    <canvas
      ref={canvasRef}
      className="fixed inset-0 pointer-events-none z-0"
      style={{ filter: 'blur(100px)' }}
    />
  );
}
