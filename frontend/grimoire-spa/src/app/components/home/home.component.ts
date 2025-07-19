import { Component, OnInit } from '@angular/core';
import { NoteService } from '../../services/note.service';
import { RouterModule } from '@angular/router';
import { Note } from '../../models/note.model';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [RouterModule, CommonModule],
  templateUrl: './home.component.html',
  styleUrl: './home.component.css'
})
export class HomeComponent implements OnInit {
  notes: Note[] = [];

  constructor(private noteService: NoteService) { }

  ngOnInit() {
    this.noteService.getNotes().subscribe(notes => {
      console.log('Fetched notes:', notes);
      this.notes = notes;
    });
  }

  deleteNote(id: string) {
    this.noteService.deleteNote(id).subscribe(() => {
      this.notes = this.notes.filter(note => note.id !== id);
    });
  }
}
