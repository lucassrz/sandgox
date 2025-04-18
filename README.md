# Sandgox

## Problématique
Créer un jeu sandbox composé de différents élements avec comme
objectif un maximum de FPS

## Le projet

### Micro optimisation
- Mémoïsation des objets rectangle de largeur n.
- On réutilise au maximum les objets créés.
- Les fonctions gérants la physique ne sont pas instanciées par les cellules de la grille, de même pour la densitée. 

### Macro optimisation
- On récupère la frame précédente et on dessine par-dessus les carrés qui sont en mouvement.
- On dessine des rectangles horizontalement si des élements de la même couleur sont collés.

### Benchmark
Comme le programme fonctionne avec une interface graphique,
nous avons du implémenter notre propre mode de manière de
benchmark le programme.

Pour cela, nous avons ajouté un flag falcultatif
"-benchmark true" au lancement du programme. Il permet de
rentrer dans ce mode et la fenêtre ce fermera lorsque le
programme aura passé la 100ème frame.

C'est notre valeur de mesure pour ce benchmark.

Afin d'être le plus proche possible de l'utilisation,
Le benchmark est réalisé avec tous les élements
(sable, eau, métal, générateur d'eau, trou noir)
ainsi que les différents états possibles (en mouvement, statique)

### Résultats

**V1**: **31.07** secondes pour 100 frames.

**Programme final**: **8.35** secondes pour 100 frames.

Résultat obtenu sur un Macbook Pro M1 Pro 16Go RAM