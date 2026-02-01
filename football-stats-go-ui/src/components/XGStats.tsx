import React, { useState, FormEvent } from 'react';
import { ApiResponse, XGStatFixture } from '../types';
import ShotMap from './ShotMap';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function XGStats() {
  const [url, setUrl] = useState<string>('');
  const [fixtureId, setFixtureId] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [fixtures, setFixtures] = useState<XGStatFixture[]>([]);
  const [currentFixtureIndex, setCurrentFixtureIndex] = useState<number>(0);
  const [error, setError] = useState<string | null>(null);

  const handleScrape = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/api/scrape/xgstats`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      });

      const data: ApiResponse<XGStatFixture> = await response.json();

      if (response.ok && data.success && data.data) {
        // Add to fixtures list and set as current
        setFixtures(prev => [data.data!, ...prev]);
        setCurrentFixtureIndex(0);
        setUrl('');
      } else {
        setError(data.error || 'Failed to scrape fixture');
      }
    } catch (err) {
      setError((err as Error).message || 'Network error');
    } finally {
      setLoading(false);
    }
  };

  const handleLoadFixture = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/api/xgstats?id=${fixtureId}`);
      const data: ApiResponse<XGStatFixture> = await response.json();

      if (response.ok && data.success && data.data) {
        const existingIndex = fixtures.findIndex(f => f.id === data.data!.id);
        if (existingIndex >= 0) {
          setCurrentFixtureIndex(existingIndex);
        } else {
          setFixtures(prev => [data.data!, ...prev]);
          setCurrentFixtureIndex(0);
        }
        setFixtureId('');
      } else {
        setError(data.error || 'Failed to load fixture');
      }
    } catch (err) {
      setError((err as Error).message || 'Network error');
    } finally {
      setLoading(false);
    }
  };

  const currentFixture = fixtures[currentFixtureIndex];

  const nextFixture = () => {
    if (currentFixtureIndex < fixtures.length - 1) {
      setCurrentFixtureIndex(currentFixtureIndex + 1);
    }
  };

  const prevFixture = () => {
    if (currentFixtureIndex > 0) {
      setCurrentFixtureIndex(currentFixtureIndex - 1);
    }
  };

  return (
    <div className="xgstats-container">
      <div className="scraper-section">
        <h2>⚽ xG Statistics</h2>
        
        <div className="input-forms">
          <form onSubmit={handleScrape} className="scrape-form">
            <div className="form-group">
              <label htmlFor="url">Scrape from xgstat.com</label>
              <div className="input-button-group">
                <input
                  id="url"
                  type="url"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  placeholder="https://xgstat.com/fixture/..."
                  required
                />
                <button type="submit" disabled={loading}>
                  {loading ? '⏳' : '🔍'} Scrape
                </button>
              </div>
            </div>
          </form>

          <form onSubmit={handleLoadFixture} className="load-form">
            <div className="form-group">
              <label htmlFor="fixtureId">Load Saved Fixture</label>
              <div className="input-button-group">
                <input
                  id="fixtureId"
                  type="number"
                  value={fixtureId}
                  onChange={(e) => setFixtureId(e.target.value)}
                  placeholder="Fixture ID"
                  required
                />
                <button type="submit" disabled={loading}>
                  {loading ? '⏳' : '📥'} Load
                </button>
              </div>
            </div>
          </form>
        </div>

        {error && (
          <div className="error-message">
            <strong>Error:</strong> {error}
          </div>
        )}
      </div>

      {currentFixture && (
        <div className="fixture-display">
          {/* Navigation */}
          {fixtures.length > 1 && (
            <div className="fixture-nav">
              <button 
                onClick={prevFixture} 
                disabled={currentFixtureIndex === 0}
                className="nav-button"
              >
                ← Previous
              </button>
              <span className="fixture-counter">
                {currentFixtureIndex + 1} / {fixtures.length}
              </span>
              <button 
                onClick={nextFixture} 
                disabled={currentFixtureIndex === fixtures.length - 1}
                className="nav-button"
              >
                Next →
              </button>
            </div>
          )}

          {/* Match Header */}
          <div className="match-header">
            <div className="team home-team">
              <h2>{currentFixture.home_team}</h2>
              <div className="score">{currentFixture.home_score}</div>
              <div className="xg-badge">xG: {currentFixture.home_xg.toFixed(2)}</div>
            </div>
            
            <div className="match-info">
              <div className="vs">VS</div>
              <div className="match-details">
                <div>Gameweek {currentFixture.gameweek}</div>
                <div>{new Date(currentFixture.date).toLocaleDateString()}</div>
                <div className="fixture-id">ID: {currentFixture.id}</div>
              </div>
            </div>

            <div className="team away-team">
              <h2>{currentFixture.away_team}</h2>
              <div className="score">{currentFixture.away_score}</div>
              <div className="xg-badge">xG: {currentFixture.away_xg.toFixed(2)}</div>
            </div>
          </div>

          {/* Shot Map */}
          <ShotMap
            homeShots={currentFixture.home_shots}
            awayShots={currentFixture.away_shots}
            homeTeam={currentFixture.home_team}
            awayTeam={currentFixture.away_team}
          />

          {/* Shot Lists */}
          <div className="shots-section">
            <div className="shots-column">
              <h3>{currentFixture.home_team} Shots ({currentFixture.home_shots.length})</h3>
              <div className="shots-list">
                {currentFixture.home_shots
                  .sort((a, b) => a.minute - b.minute)
                  .map((shot, idx) => (
                    <div key={idx} className={`shot-item ${shot.is_goal ? 'goal' : ''}`}>
                      <div className="shot-header">
                        <span className="minute">{shot.minute}'</span>
                        <span className="player">{shot.player_name}</span>
                        {shot.is_goal && <span className="goal-icon">⚽</span>}
                      </div>
                      <div className="shot-details">
                        <span className="xg-value">xG: {shot.xg.toFixed(2)}</span>
                        <span className="shot-type">{shot.shot_type}</span>
                      </div>
                    </div>
                  ))}
              </div>
            </div>

            <div className="shots-column">
              <h3>{currentFixture.away_team} Shots ({currentFixture.away_shots.length})</h3>
              <div className="shots-list">
                {currentFixture.away_shots
                  .sort((a, b) => a.minute - b.minute)
                  .map((shot, idx) => (
                    <div key={idx} className={`shot-item ${shot.is_goal ? 'goal' : ''}`}>
                      <div className="shot-header">
                        <span className="minute">{shot.minute}'</span>
                        <span className="player">{shot.player_name}</span>
                        {shot.is_goal && <span className="goal-icon">⚽</span>}
                      </div>
                      <div className="shot-details">
                        <span className="xg-value">xG: {shot.xg.toFixed(2)}</span>
                        <span className="shot-type">{shot.shot_type}</span>
                      </div>
                    </div>
                  ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {!currentFixture && !loading && (
        <div className="empty-state">
          <p>🎯 Scrape a fixture from xgstat.com or load a saved fixture to view shot maps and statistics</p>
        </div>
      )}
    </div>
  );
}

export default XGStats;
