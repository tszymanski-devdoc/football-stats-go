import React, { useState, FormEvent } from 'react';
import { MatchData, MatchPrediction, ApiResponse } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function PredictMatch() {
  const [homeTeam, setHomeTeam] = useState<string>('');
  const [awayTeam, setAwayTeam] = useState<string>('');
  const [homeMatchesJson, setHomeMatchesJson] = useState<string>('');
  const [awayMatchesJson, setAwayMatchesJson] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [result, setResult] = useState<MatchPrediction | null>(null);
  const [error, setError] = useState<string | null>(null);

  const sampleHomeData = `[
  {
    "id": "1",
    "home_team": "Manchester United",
    "away_team": "Chelsea",
    "home_score": 2,
    "away_score": 1,
    "match_date": "2024-01-15",
    "league": "Premier League",
    "season": "2023-24"
  }
]`;

  const sampleAwayData = `[
  {
    "id": "2",
    "home_team": "Arsenal",
    "away_team": "Liverpool",
    "home_score": 1,
    "away_score": 2,
    "match_date": "2024-01-20",
    "league": "Premier League",
    "season": "2023-24"
  }
]`;

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const homeMatches: MatchData[] = JSON.parse(homeMatchesJson);
      const awayMatches: MatchData[] = JSON.parse(awayMatchesJson);

      const response = await fetch(`${API_URL}/api/predict-match`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          home_team: homeTeam,
          away_team: awayTeam,
          home_matches: homeMatches,
          away_matches: awayMatches,
        }),
      });

      const data: ApiResponse<MatchPrediction> = await response.json();

      if (response.ok && data.success) {
        setResult(data.data!);
      } else {
        setError(data.error || 'Failed to predict match');
      }
    } catch (err) {
      setError((err as Error).message || 'Invalid JSON or network error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-container">
      <h2>🎯 Predict Match Outcome</h2>
      <form onSubmit={handleSubmit}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
          <div className="form-group">
            <label htmlFor="homeTeam">Home Team</label>
            <input
              id="homeTeam"
              type="text"
              value={homeTeam}
              onChange={(e) => setHomeTeam(e.target.value)}
              placeholder="e.g., Manchester United"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="awayTeam">Away Team</label>
            <input
              id="awayTeam"
              type="text"
              value={awayTeam}
              onChange={(e) => setAwayTeam(e.target.value)}
              placeholder="e.g., Liverpool"
              required
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="homeMatches">Home Team Matches (JSON)</label>
          <textarea
            id="homeMatches"
            value={homeMatchesJson}
            onChange={(e) => setHomeMatchesJson(e.target.value)}
            placeholder={sampleHomeData}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="awayMatches">Away Team Matches (JSON)</label>
          <textarea
            id="awayMatches"
            value={awayMatchesJson}
            onChange={(e) => setAwayMatchesJson(e.target.value)}
            placeholder={sampleAwayData}
            required
          />
          <small>Enter historical match data for both teams as JSON arrays.</small>
        </div>

        <div style={{ display: 'flex', gap: '1rem' }}>
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Predicting...' : 'Predict Match'}
          </button>
          <button
            type="button"
            className="btn-primary"
            onClick={() => {
              setHomeMatchesJson(sampleHomeData);
              setAwayMatchesJson(sampleAwayData);
            }}
            style={{ background: '#6c757d' }}
          >
            Load Sample
          </button>
        </div>
      </form>

      {loading && <div className="loading">Predicting match outcome...</div>}

      {error && (
        <div className="result-container error">
          <h3>Error</h3>
          <pre>{error}</pre>
        </div>
      )}

      {result && (
        <div className="result-container">
          <h3>Match Prediction</h3>
          <div style={{ marginTop: '1rem' }}>
            <div style={{ background: 'white', padding: '1.5rem', borderRadius: '8px', marginBottom: '1rem', textAlign: 'center' }}>
              <h4 style={{ marginBottom: '0.5rem', color: '#667eea' }}>
                {result.home_team} vs {result.away_team}
              </h4>
              <div style={{ fontSize: '2rem', fontWeight: 'bold', margin: '1rem 0' }}>
                {result.predicted_score}
              </div>
              <div style={{ fontSize: '0.9rem', color: '#666' }}>
                Confidence: {(result.confidence * 100).toFixed(1)}%
              </div>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem' }}>
              <div style={{ 
                background: 'white', 
                padding: '1rem', 
                borderRadius: '8px', 
                textAlign: 'center',
                border: result.home_win_probability > result.away_win_probability && result.home_win_probability > result.draw_probability ? '2px solid #667eea' : 'none'
              }}>
                <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.5rem' }}>Home Win</div>
                <div style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#667eea' }}>
                  {(result.home_win_probability * 100).toFixed(1)}%
                </div>
              </div>
              <div style={{ 
                background: 'white', 
                padding: '1rem', 
                borderRadius: '8px', 
                textAlign: 'center',
                border: result.draw_probability > result.home_win_probability && result.draw_probability > result.away_win_probability ? '2px solid #667eea' : 'none'
              }}>
                <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.5rem' }}>Draw</div>
                <div style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#6c757d' }}>
                  {(result.draw_probability * 100).toFixed(1)}%
                </div>
              </div>
              <div style={{ 
                background: 'white', 
                padding: '1rem', 
                borderRadius: '8px', 
                textAlign: 'center',
                border: result.away_win_probability > result.home_win_probability && result.away_win_probability > result.draw_probability ? '2px solid #667eea' : 'none'
              }}>
                <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.5rem' }}>Away Win</div>
                <div style={{ fontSize: '1.5rem', fontWeight: 'bold', color: '#764ba2' }}>
                  {(result.away_win_probability * 100).toFixed(1)}%
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default PredictMatch;
