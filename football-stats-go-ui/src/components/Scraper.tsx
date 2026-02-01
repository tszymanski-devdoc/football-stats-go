import React, { useState, FormEvent } from 'react';
import { ApiResponse } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

interface ScrapeData {
  url: string;
  title: string;
  content: string;
  timestamp: string;
  metadata: {
    scraped_at: string;
    status: string;
  };
}

function Scraper() {
  const [url, setUrl] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [result, setResult] = useState<ScrapeData | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await fetch(`${API_URL}/api/scrape`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url }),
      });

      const data: ApiResponse<ScrapeData> = await response.json();

      if (response.ok && data.success) {
        setResult(data.data!);
      } else {
        setError(data.error || 'Failed to scrape website');
      }
    } catch (err) {
      setError((err as Error).message || 'Network error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-container">
      <h2>🌐 Scrape Website</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="url">Website URL</label>
          <input
            id="url"
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://example.com"
            required
          />
          <small>Enter the URL of the website you want to scrape.</small>
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Scraping...' : 'Scrape Website'}
        </button>
      </form>

      {error && (
        <div className="error-message">
          <strong>Error:</strong> {error}
        </div>
      )}

      {result && (
        <div className="result-container">
          <h3>Scrape Results</h3>
          <div className="result-card">
            <div className="result-row">
              <strong>URL:</strong>
              <span>{result.url}</span>
            </div>
            <div className="result-row">
              <strong>Title:</strong>
              <span>{result.title}</span>
            </div>
            <div className="result-row">
              <strong>Scraped At:</strong>
              <span>{new Date(result.metadata.scraped_at).toLocaleString()}</span>
            </div>
            <div className="result-row">
              <strong>Status:</strong>
              <span className="badge success">{result.metadata.status}</span>
            </div>
          </div>

          <div className="content-section">
            <h4>Page Content Preview</h4>
            <div className="content-preview">
              {result.content.substring(0, 500)}
              {result.content.length > 500 && '...'}
            </div>
            <small className="text-muted">
              Showing first 500 characters of {result.content.length} total characters
            </small>
          </div>
        </div>
      )}
    </div>
  );
}

export default Scraper;
