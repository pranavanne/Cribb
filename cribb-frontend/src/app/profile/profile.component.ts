import { Component } from "@angular/core";

@Component({
    selector: 'app-profile',
    templateUrl: './profile.component.html',
})

export class ProfileComponent {
    user = {
        name: 'John Doe',
        email: 'johndoe@example.com',
        phone: '+1234567890',
        group: 'Sunset Apartments',
        room: 'A-202',
        cribbPoints: 7
    };

    getRanking(points: number): string {
        if (points > 5) return 'Good';
        if (points < 5) return 'Poor';
        return 'Neutral';
    }
}