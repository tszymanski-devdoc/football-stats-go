import React, { useState, FormEvent } from 'react';
import { MatchData, TeamStats, ApiResponse } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

function AnalyzeTeam() {
  const [teamName, setTeamName] = useState<string>('');
  const [matchesJson, setMatchesJson] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [result, setResult] = useState<TeamStats | null>(null);
  const [error, setError] = useState<string | null>(null);

  const sampleData = `[
  {
    "id": "1",
    "home_team": "Manchester United",
    "away_team": "Chelsea",
    "home_score": 2,
    "away_score": 1,
    "match_date": "2024-01-15T15:00:00Z",
    "league": "Premier League",
    "season": "2023-24"
  },
  {
    "id": "2",
    "home_team": "Liverpool",
    "away_team": "Manchester United",
    "home_score": 1,
    "away_score": 2,
    "match_date": "2024-01-22T15:00:00Z",
    "league": "Premier League",
    "season": "2023-24"
  }
]`;
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
      const matches: MatchData[] = JSON.parse(matchesJson);
      const response = await fetch(`${API_URL}/api/analyze-team`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          team_name: teamName,
          matches: matches,
        }),
      });

      const data: ApiResponse<TeamStats> = await response.json();

      if (response.ok && data.success) {
        setResult(data.data!);
      } else {
        setError(data.error || 'Failed to analyze team');
      }
    } catch (err) {
      setError((err as Error).message || 'Invalid JSON or network error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-container">
      <h2>📊 Analyze Team Statistics</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="teamName">Team Name</label>
          <input
            id="teamName"
            type="text"
            value={teamName}
            onChange={(e) => setTeamName(e.target.value)}
            placeholder="e.g., Manchester United"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="matches">Match Data (JSON)</label>
          <textarea
            id="matches"
            value={matchesJson}
            onChange={(e) => setMatchesJson(e.target.value)}
            placeholder={sampleData}
            required
          />
          <small>Enter match data as JSON array. Click "Load Sample" to see format.</small>
        </div>

        <div style={{ display: 'flex', gap: '1rem' }}>
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Analyzing...' : 'Analyze Team'}
          </button>
          <button
            type="button"
            className="btn-primary"
            onClick={() => setMatchesJson(sampleData)}
            style={{ background: '#6c757d' }}
          >
            Load Sample
          </button>
        </div>
      </form>

      {loading && <div className="loading">Analyzing team statistics...</div>}

      {error && (
        <div className="result-container error">
          <h3>Error</h3>
          <pre>{error}</pre>
        </div>
      )}

      {result && (
        <div className="result-container">
          <h3>Team Statistics</h3>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem', marginTop: '1rem' }}>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Team:</strong> {result.team_name}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Matches:</strong> {result.matches_played}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Wins:</strong> {result.wins}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Draws:</strong> {result.draws}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Losses:</strong> {result.losses}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Win %:</strong> {(result.win_percentage * 100).toFixed(1)}%
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Goals For:</strong> {result.goals_for}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Goals Against:</strong> {result.goals_against}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Avg Scored:</strong> {result.avg_goals_scored.toFixed(2)}
            </div>
            <div style={{ background: 'white', padding: '1rem', borderRadius: '8px' }}>
              <strong>Avg Conceded:</strong> {result.avg_goals_conceded.toFixed(2)}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default AnalyzeTeam;
