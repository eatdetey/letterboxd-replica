from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("reviews", "0001_initial"),
    ]

    operations = [
        migrations.AddField(
            model_name="review",
            name="username",
            field=models.CharField(
                blank=True,
                default="",
                max_length=255,
                verbose_name="Имя пользователя",
            ),
        ),
    ]
