export interface XGStatShot {
  x: number;
  y: number;
  xg: number;
  is_goal: boolean;
  shot_type: string;
  player_name: string;
  minute: number;
}

export interface XGStatFixture {
  gameweek: number;
  id: number;
  date: string;
  home_team: string;
  away_team: string;
  home_score: number;
  away_score: number;
  home_xg: number;
  away_xg: number;
  home_shots: XGStatShot[];
  away_shots: XGStatShot[];
}

export interface MatchData {
  id: string;
  home_team: string;
  away_team: string;
  home_score: number;
  away_score: number;
  match_date: string;
  league: string;
  season: string;
}

export interface TeamStats {
  team_name: string;
  matches_played: number;
  wins: number;
  draws: number;
  losses: number;
  goals_for: number;
  goals_against: number;
  win_percentage: number;
  avg_goals_scored: number;
  avg_goals_conceded: number;
}

export interface MatchPrediction {
  home_team: string;
  away_team: string;
  predicted_score: string;
  home_win_probability: number;
  draw_probability: number;
  away_win_probability: number;
  confidence: number;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}
