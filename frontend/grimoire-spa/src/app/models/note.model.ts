export interface Note {
    id: string;
    title: string;
    content: string;
    created_at: Date;
    updated_at: Date;
    user_id: string;
    is_public: boolean;
    tags: Tag[];
}

export interface Tag {
    id: string;
    name: string;
}
