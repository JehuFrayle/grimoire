export interface User {
    id: string;
    username: string;
    email: string;
    created_at: Date;
    updated_at: Date;
    last_login: Date;
    role: string;
    active: boolean;
    profile: Profile;
}

export interface Profile {
    first_name: string;
    last_name: string;
    bio: string;
    avatar_url: string;
    links: Link[];
}

export interface Link {
    url: string;
    title: string;
    icon: string;
    active: boolean;
}
