from django.db import models


class Genre(models.Model):
    name = models.CharField(max_length=15)


class Movie(models.Model):
    title = models.CharField(max_length=255)
    original_title = models.CharField(max_length=255, null=True)
    description = models.CharField(max_length=500, null=True)
    release_date = models.DateField(null=True)
    duration_minutes = models.IntegerField(null=True)
    country = models.CharField(max_length=20, null=True)
    age_rating = models.CharField(max_length=10, null=True)
    genres = models.ManyToManyField(Genre)    


    def __str__(self):
        return f"{self.title} ({self.release_date.year})"


class Person(models.Model):
    name = models.CharField(max_length=80)
    birth_date = models.DateField(null=True)
    biography = models.CharField(max_length=500)


class MovieCast(models.Model):
    movie = models.ForeignKey(Movie, on_delete=models.CASCADE)
    person = models.ForeignKey(Person, on_delete=models.CASCADE)
    role_type = models.CharField(max_length=20)
    character_name = models.CharField(max_length=50)
