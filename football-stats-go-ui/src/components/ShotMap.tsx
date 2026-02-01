import React from 'react';
import { XGStatShot } from '../types';

interface ShotMapProps {
  homeShots: XGStatShot[];
  awayShots: XGStatShot[];
  homeTeam: string;
  awayTeam: string;
}

const ShotMap: React.FC<ShotMapProps> = ({ homeShots, awayShots, homeTeam, awayTeam }) => {
  const pitchWidth = 600;
  const pitchHeight = 400;

  const renderShot = (shot: XGStatShot, isHome: boolean, index: number) => {
    // Convert coordinates to SVG coordinates
    const x = (shot.x / 100) * pitchWidth;
    const y = (shot.y / 100) * pitchHeight;
    
    // Size based on xG value
    const radius = 4 + (shot.xg * 12);
    
    // Color based on outcome and team
    const fillColor = shot.is_goal 
      ? (isHome ? '#10b981' : '#f59e0b')  // Green for home goals, amber for away goals
      : (isHome ? '#3b82f6' : '#ef4444'); // Blue for home shots, red for away shots
    
    const opacity = 0.3 + (shot.xg * 0.6);

    return (
      <g key={`${isHome ? 'home' : 'away'}-${index}`}>
        <circle
          cx={x}
          cy={y}
          r={radius}
          fill={fillColor}
          opacity={opacity}
          stroke={shot.is_goal ? '#fff' : 'none'}
          strokeWidth={shot.is_goal ? 2 : 0}
        >
          <title>
            {shot.player_name} - {shot.minute}' - xG: {shot.xg.toFixed(2)}
            {shot.is_goal ? ' ⚽ GOAL' : ''} - {shot.shot_type}
          </title>
        </circle>
        {shot.is_goal && (
          <text
            x={x}
            y={y}
            textAnchor="middle"
            dominantBaseline="middle"
            fill="#fff"
            fontSize="10"
            fontWeight="bold"
          >
            ⚽
          </text>
        )}
      </g>
    );
  };

  return (
    <div className="shot-map">
      <div className="shot-map-header">
        <div className="team-label home">{homeTeam}</div>
        <div className="shot-map-title">Shot Map</div>
        <div className="team-label away">{awayTeam}</div>
      </div>
      
      <svg 
        viewBox={`0 0 ${pitchWidth} ${pitchHeight}`}
        className="pitch-svg"
        style={{ maxWidth: '100%', height: 'auto' }}
      >
        {/* Pitch background */}
        <rect width={pitchWidth} height={pitchHeight} fill="#1a472a" />
        
        {/* Center circle */}
        <circle 
          cx={pitchWidth / 2} 
          cy={pitchHeight / 2} 
          r={60} 
          fill="none" 
          stroke="rgba(255,255,255,0.3)" 
          strokeWidth="2" 
        />
        <circle 
          cx={pitchWidth / 2} 
          cy={pitchHeight / 2} 
          r={3} 
          fill="rgba(255,255,255,0.5)" 
        />
        
        {/* Halfway line */}
        <line 
          x1={pitchWidth / 2} 
          y1="0" 
          x2={pitchWidth / 2} 
          y2={pitchHeight} 
          stroke="rgba(255,255,255,0.3)" 
          strokeWidth="2" 
        />
        
        {/* Penalty areas */}
        <rect 
          x="0" 
          y={(pitchHeight - 200) / 2} 
          width="120" 
          height="200" 
          fill="none" 
          stroke="rgba(255,255,255,0.3)" 
          strokeWidth="2" 
        />
        <rect 
          x={pitchWidth - 120} 
          y={(pitchHeight - 200) / 2} 
          width="120" 
          height="200" 
          fill="none" 
          stroke="rgba(255,255,255,0.3)" 
          strokeWidth="2" 
        />
        
        {/* Goals */}
        <rect x="0" y={(pitchHeight - 80) / 2} width="5" height="80" fill="rgba(255,255,255,0.5)" />
        <rect x={pitchWidth - 5} y={(pitchHeight - 80) / 2} width="5" height="80" fill="rgba(255,255,255,0.5)" />
        
        {/* Render shots */}
        {homeShots.map((shot, index) => renderShot(shot, true, index))}
        {awayShots.map((shot, index) => renderShot(shot, false, index))}
      </svg>

      <div className="shot-map-legend">
        <div className="legend-item">
          <div className="legend-circle" style={{ backgroundColor: '#3b82f6', opacity: 0.7 }}></div>
          <span>{homeTeam} Shot</span>
        </div>
        <div className="legend-item">
          <div className="legend-circle" style={{ backgroundColor: '#10b981', opacity: 0.9 }}></div>
          <span>{homeTeam} Goal</span>
        </div>
        <div className="legend-item">
          <div className="legend-circle" style={{ backgroundColor: '#ef4444', opacity: 0.7 }}></div>
          <span>{awayTeam} Shot</span>
        </div>
        <div className="legend-item">
          <div className="legend-circle" style={{ backgroundColor: '#f59e0b', opacity: 0.9 }}></div>
          <span>{awayTeam} Goal</span>
        </div>
        <div className="legend-item">
          <span style={{ fontSize: '12px', color: '#9ca3af' }}>Size = xG value</span>
        </div>
      </div>
    </div>
  );
};

export default ShotMap;
