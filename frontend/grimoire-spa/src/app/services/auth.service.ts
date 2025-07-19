import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor() { }

  public isAuthenticated(): boolean {
    const token = localStorage.getItem('session_token');
    return !!token;
  }

  public login(token: string): void {
    localStorage.setItem('session_token', token);
  }

  public logout(): void {
    localStorage.removeItem('session_token');
  }
}