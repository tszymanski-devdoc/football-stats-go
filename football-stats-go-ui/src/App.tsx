import React, { useState } from 'react';
import './App.css';
import AnalyzeTeam from './components/AnalyzeTeam';
import PredictMatch from './components/PredictMatch';

type ViewType = 'home' | 'analyze' | 'predict';

function App() {
  const [activeView, setActiveView] = useState<ViewType>('home');

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
            className={activeView === 'analyze' ? 'active' : ''} 
            onClick={() => setActiveView('analyze')}
          >
            Analyze Team
          </button>
          <button 
            className={activeView === 'predict' ? 'active' : ''} 
            onClick={() => setActiveView('predict')}
          >
            Predict Match
          </button>
        </div>
      </nav>

      <main className="main-content">
        {activeView === 'home' && (
          <div className="home">
            <h1>Welcome to Football Stats</h1>
            <p>Analyze team performance and predict match outcomes using advanced statistics.</p>
            <div className="features">
              <div className="feature-card">
                <h3>📊 Team Analysis</h3>
                <p>Get comprehensive statistics including win rate, goals, and performance metrics.</p>
              </div>
              <div className="feature-card">
                <h3>🎯 Match Prediction</h3>
                <p>Predict match outcomes using Poisson distribution and historical data.</p>
              </div>
            </div>
          </div>
        )}
        {activeView === 'analyze' && <AnalyzeTeam />}
        {activeView === 'predict' && <PredictMatch />}
      </main>
    </div>
  );
}

export default App;

