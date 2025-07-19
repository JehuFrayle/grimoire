import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Note } from '../models/note.model';

@Injectable({
  providedIn: 'root'
})
export class NoteService {
  private apiUrl = 'http://localhost:8080/api/notes';

  constructor(private http: HttpClient) { }

  private getAuthHeaders() {
    const token = localStorage.getItem('session_token');
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }

  getNotes(): Observable<Note[]> {
    return this.http.get<Note[]>(this.apiUrl, { headers: this.getAuthHeaders() });
  }

  getNote(id: string): Observable<Note> {
    return this.http.get<Note>(`${this.apiUrl}/${id}`, { headers: this.getAuthHeaders() });
  }

  createNote(note: Note): Observable<Note> {
    return this.http.post<Note>(this.apiUrl, note, { headers: this.getAuthHeaders() });
  }

  updateNote(id: string, note: Note): Observable<Note> {
    return this.http.patch<Note>(`${this.apiUrl}/${id}`, note, { headers: this.getAuthHeaders() });
  }

  deleteNote(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`, { headers: this.getAuthHeaders() });
  }
}