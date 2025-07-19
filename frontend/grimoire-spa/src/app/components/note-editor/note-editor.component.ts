import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NoteService } from '../../services/note.service';
import { Note } from '../../models/note.model';

@Component({
  selector: 'app-note-editor',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './note-editor.component.html',
  styleUrl: './note-editor.component.css'
})
export class NoteEditorComponent implements OnInit {
  noteForm: FormGroup;
  noteId: string | null = null;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private noteService: NoteService
  ) {
    this.noteForm = this.fb.group({
      title: ['', Validators.required],
      content: ['', Validators.required]
    });
  }

  ngOnInit() {
    this.noteId = this.route.snapshot.paramMap.get('id');
    if (this.noteId) {
      this.noteService.getNote(this.noteId).subscribe((note: Note) => {
        this.noteForm.patchValue(note);
      });
    }
  }

  onSubmit() {
    if (this.noteForm.valid) {
      if (this.noteId) {
        this.noteService.updateNote(this.noteId, this.noteForm.value).subscribe(() => {
          this.router.navigate(['/']);
        });
      } else {
        this.noteService.createNote(this.noteForm.value).subscribe(() => {
          this.router.navigate(['/']);
        });
      }
    }
  }
}
