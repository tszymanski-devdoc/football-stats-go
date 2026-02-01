import React, { useState } from 'react';
import './App.css';
import XGStats from './components/XGStats';

type ViewType = 'home' | 'xgstats';

function App() {
  const [activeView, setActiveView] = useState<ViewType>('xgstats');

  return (
    <div className="App">
      <nav className="navbar">
        <div className="nav-brand">⚽ Football Stats</div>
        <div className="nav-links">
          <button 
            className={activeView === 'home' ? 'active' : ''} 
            onClick={() => setActiveView('home')}
          >
            Home
          </button>
          <button 
            className={activeView === 'xgstats' ? 'active' : ''} 
            onClick={() => setActiveView('xgstats')}
          >
            xG Stats
          </button>
        </div>
      </nav>

      <main className="main-content">
        {activeView === 'home' && (
          <div className="home">
            <h1>Welcome to Football Stats</h1>
            <p>View and analyze xG (expected goals) statistics from Premier League matches.</p>
            <div className="features">
              <div className="feature-card" onClick={() => setActiveView('xgstats')} style={{ cursor: 'pointer' }}>
                <h3>📊 xG Statistics</h3>
                <p>Interactive shot maps and xG analysis for football matches.</p>
              </div>
            </div>
          </div>
        )}
        {activeView === 'xgstats' && <XGStats />}
      </main>
    </div>
  );
}

export default App;